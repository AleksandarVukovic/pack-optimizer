package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/aleksandarv/pack-optimizer/gen/http/optimizer/server"
	goaoptimizer "github.com/aleksandarv/pack-optimizer/gen/optimizer"
	"github.com/aleksandarv/pack-optimizer/internal/api"
	"github.com/aleksandarv/pack-optimizer/internal/calculator"
	"github.com/aleksandarv/pack-optimizer/internal/logger"
	"github.com/aleksandarv/pack-optimizer/internal/pack"
	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/http/middleware"
)

func main() {
	var (
		debug bool
		port  int
	)
	flag.BoolVar(&debug, "debug", false, "Enable debug mode with verbose logging")
	flag.IntVar(&port, "port", 8080, "HTTP port")
	loadFlagsFromEnv()
	flag.Parse()

	log := logger.NewLogger(debug)

	var (
		ctx context.Context
	)
	{
		ctx = logger.WithCtx(context.Background(), log)
	}

	psvc := pack.NewInMemorySvc(pack.DefaultSizes)
	calculator := calculator.NewCalculator(psvc)
	optimizerSvc := api.NewOptimizerSvc(psvc, calculator)
	endpoints := goaoptimizer.NewEndpoints(optimizerSvc)

	mux := goahttp.NewMuxer()
	optimizerSrv := server.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, errorHandler(log), nil)

	optimizerSrv.Use(logger.RequestMiddleware(log))
	optimizerSrv.Use(middleware.PopulateRequestContext())
	optimizerSrv.Use(middleware.RequestID(
		middleware.UseXRequestIDHeaderOption(true),
		middleware.XRequestHeaderLimitOption(64),
	))

	server.Mount(mux, optimizerSrv)
	for _, m := range optimizerSrv.Mounts {
		log.Debug("expose API", "verb", m.Verb, "path", m.Pattern, "method", m.Method)
	}

	addr := ":" + strconv.Itoa(port)
	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: time.Second * 60}

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		log.Info("Start server on", "host", addr)
		errc <- srv.ListenAndServe()
	}()

	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Go(func() {
		<-ctx.Done()
		log.Info("shutting down server")

		// do shutdown with 30s timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("error while shutting down", "error", err)
		}
	})

	// waiting on some signal to shutdown the server
	err := <-errc
	log.Info("exiting server", "reason", err)

	// trigger shutdown goroutine process
	cancel()
	wg.Wait()
}

func loadFlagsFromEnv() {
	envToFlag := map[string]string{
		"DEBUG": "debug",
		"PORT":  "port",
	}
	for env, flagName := range envToFlag {
		if val := os.Getenv(env); val != "" {
			os.Args = append(os.Args, fmt.Sprintf("--%s=%s", flagName, val))
		}
	}
}

func errorHandler(log *slog.Logger) func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		if verr, ok := errors.AsType[*pack.ValidationError](err); ok {
			log.Error("validation error", "error", verr.Error())
			return
		}
		log.Error("GOA error", "error", err.Error())
	}
}
