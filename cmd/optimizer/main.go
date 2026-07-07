package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
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
	"github.com/aleksandarv/pack-optimizer/internal/telemetry"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/http/middleware"
)

func main() {
	var (
		debug           bool
		port            int
		metricsInterval int
		otelEndpoint    string
		dbHost          string
		dbPort          int
		dbUser          string
		dbPassword      string
	)
	flag.BoolVar(&debug, "debug", false, "Enable debug mode with verbose logging")
	flag.IntVar(&port, "port", 8080, "HTTP port")
	flag.IntVar(&metricsInterval, "metrics-interval", 60, "Interval in seconds for exporting metrics")
	flag.StringVar(&otelEndpoint, "otel-endpoint", "localhost:4318", "OTel collector endpoint")
	flag.StringVar(&dbHost, "db-host", "localhost", "PostgreSQL host")
	flag.IntVar(&dbPort, "db-port", 5432, "PostgreSQL port")
	flag.StringVar(&dbUser, "db-user", "postgres", "PostgreSQL user")
	flag.StringVar(&dbPassword, "db-password", "postgres", "PostgreSQL password")
	loadFlagsFromEnv()
	flag.Parse()

	log := logger.NewLogger(debug)
	ctx := logger.WithCtx(context.Background(), log)

	dbpool, err := pgxpool.New(ctx, buildDbUrl(dbUser, dbPassword, dbHost, dbPort))
	if err != nil {
		log.Error("unable to create postgres connection pool", "error", err)
		os.Exit(1)
	}
	if err := dbpool.Ping(ctx); err != nil {
		log.Error("unable to reach postgres", "error", err)
		os.Exit(1)
	}
	log.Info("connected to postgres")

	psvc := pack.NewPersistentSvc(dbpool, telemetry.NewTracer(telemetry.PackComponentName))
	calculator := calculator.NewCalculator(psvc)
	meter := telemetry.NewMeter(telemetry.APIComponentName)
	optimizerSvc := api.NewOptimizerSvc(psvc, calculator, meter)
	endpoints := goaoptimizer.NewEndpoints(optimizerSvc)

	mux := goahttp.NewMuxer()
	optimizerSrv := server.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, errorHandler(log), nil)

	optimizerSrv.Use(telemetry.RequestIdMiddleware())
	optimizerSrv.Use(otelhttp.NewMiddleware(telemetry.APIComponentName))
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

	// temporary solution just to show index page and docs
	mountDocsEndpoints(ctx, mux)

	shutdownTracing, shutdownMetrics, err := telemetry.InitObservability(ctx, otelEndpoint, metricsInterval)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
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
		if err := shutdownTracing(ctx); err != nil {
			log.Error("error while shutting down tracer provider", "error", err)
		}
		if err := shutdownMetrics(ctx); err != nil {
			log.Error("error while shutting down meter provider", "error", err)
		}
		dbpool.Close()
	})

	// waiting on some signal to shutdown the server
	err = <-errc
	log.Info("exiting server", "reason", err)

	// trigger shutdown goroutine process
	cancel()
	wg.Wait()
}

// mount few APIs about docs and index page
func mountDocsEndpoints(ctx context.Context, mux goahttp.ResolverMuxer) {
	// docker build will ensure that needed files are shipped together with binary
	log := logger.FromCtx(ctx)
	mux.Handle("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, filepath.Join("index.html"))
	})
	log.Debug("expose API", "verb", "GET", "path", "/")

	mux.Handle("GET", "/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, filepath.Join("openapi3.json"))
	})
	log.Debug("expose API", "verb", "GET", "path", "/openapi.json")

	mux.Handle("GET", "/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, filepath.Join("swagger.html"))
	})
	log.Debug("expose API", "verb", "GET", "path", "/docs")
}

func buildDbUrl(user, password, host string, port int) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   host + ":" + strconv.Itoa(port),
		Path:   "pack_optimizer",
	}
	return u.String() + "?sslmode=disable"
}

func loadFlagsFromEnv() {
	envToFlag := map[string]string{
		"DEBUG":            "debug",
		"PORT":             "port",
		"METRICS_INTERVAL": "metrics-interval",
		"OTEL_ENDPOINT":    "otel-endpoint",
		"DB_HOST":          "db-host",
		"DB_PORT":          "db-port",
		"DB_USER":          "db-user",
		"DB_PASSWORD":      "db-password",
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
