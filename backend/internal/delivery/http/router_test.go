package http

import (
	"io"
	"log/slog"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prayogopangestu/crm-system/backend/pkg/auth"
)

func TestHealthAndReadiness(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewHandler(nil, logger)
	router := Router(
		handler, auth.New("01234567890123456789012345678901", time.Hour),
		nil, logger, []string{"http://localhost:3000"}, func() error { return nil },
	)
	for _, path := range []string{"/healthz", "/readyz"} {
		request := httptest.NewRequest(nethttp.MethodGet, path, nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != nethttp.StatusOK {
			t.Fatalf("%s returned %d", path, response.Code)
		}
	}
}

func TestProtectedRouteRequiresBearerToken(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := NewHandler(nil, logger)
	router := Router(
		handler, auth.New("01234567890123456789012345678901", time.Hour),
		nil, logger, nil, func() error { return nil },
	)
	request := httptest.NewRequest(nethttp.MethodGet, "/api/profile", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != nethttp.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}
}
