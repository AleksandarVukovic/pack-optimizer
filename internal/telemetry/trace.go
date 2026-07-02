package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type otelTracer struct {
	tracer trace.Tracer
}

func NewTracer(name string) Tracer {
	return &otelTracer{tracer: otel.Tracer(name)}
}

func (t *otelTracer) Trace(ctx context.Context, method string, attrs ...Attribute) (context.Context, func(err *error)) {
	ctx, span := t.tracer.Start(ctx, method)
	for _, a := range attrs {
		span.SetAttributes(attribute.String(a.Key, a.Value))
	}

	return ctx, func(err *error) {
		if err != nil && *err != nil {
			span.RecordError(*err)
			span.SetStatus(codes.Error, (*err).Error())
		}
		span.End()
	}
}
