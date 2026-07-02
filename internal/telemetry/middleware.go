package telemetry

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"goa.design/goa/v3/middleware"
)

func RequestIdMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
				span := trace.SpanFromContext(r.Context())
				span.SetAttributes(attribute.String("reqID", reqID.(string)))
			}
			next.ServeHTTP(w, r)
		})
	}
}
