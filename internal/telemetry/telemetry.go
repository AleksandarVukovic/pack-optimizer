package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/aleksandarv/pack-optimizer/internal/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

var otelServiceName = semconv.ServiceNameKey.String(ServiceName)

func InitObservability(ctx context.Context, endpoint string, metricIntervalSeconds int) (shutdownTracing func(context.Context) error, shutdownMetrics func(context.Context) error, err error) {
	log := logger.FromCtx(ctx)
	log.Debug("initializing observability", "otel_endpoint", endpoint, "metrics_interval_seconds", metricIntervalSeconds)

	r, err := resource.New(ctx, resource.WithAttributes(otelServiceName))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create otel resource: %w", err)
	}

	tp, err := createTracer(ctx, r, endpoint)
	if err != nil {
		return nil, nil, err
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	mp, err := createMeter(ctx, r, endpoint, metricIntervalSeconds)
	if err != nil {
		return nil, nil, err
	}
	otel.SetMeterProvider(mp)

	log.Info("observability initialized")
	return tp.Shutdown, mp.Shutdown, nil
}

func createTracer(ctx context.Context, r *resource.Resource, endpoint string) (*trace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create otel http tracer exporter: %w", err)
	}

	return trace.NewTracerProvider(
		trace.WithResource(r),
		trace.WithBatcher(exporter),
		trace.WithSampler(trace.AlwaysSample()),
	), nil
}

func createMeter(ctx context.Context, r *resource.Resource, endpoint string, intervalSeconds int) (*metric.MeterProvider, error) {
	me, err := otlpmetrichttp.New(
		ctx,
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create otel http metric exporter: %w", err)
	}

	return metric.NewMeterProvider(
		metric.WithResource(r),
		metric.WithReader(metric.NewPeriodicReader(me, metric.WithInterval(time.Duration(intervalSeconds)*time.Second))),
	), nil
}
