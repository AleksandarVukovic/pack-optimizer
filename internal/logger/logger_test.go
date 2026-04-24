package logger

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"goa.design/goa/v3/middleware"
)

func TestWithCtxAndFromCtx(t *testing.T) {
	log := NewLogger(false)
	ctx := WithCtx(context.Background(), log)

	got := FromCtx(ctx)
	if got != log {
		t.Errorf("expected same logger, got different instance")
	}
}

func TestFromCtx_PanicWithoutLogger(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic when no logger in context")
		}
	}()
	_ = FromCtx(context.Background())
}

func TestRequestMiddleware_WithRequestID(t *testing.T) {
	log := NewLogger(false)
	mw := RequestMiddleware(log)

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromCtx(r.Context())
		if l == nil {
			t.Error("expected logger in context")
		}
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "12345")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
}

func TestRequestMiddleware_WithoutRequestID(t *testing.T) {
	log := NewLogger(false)
	mw := RequestMiddleware(log)

	// Create a handler that should not panic
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logger may or may not be in context, but should not panic
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
}
