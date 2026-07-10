package httpserver

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prayogopangestu/crm-system/backend/pkg/auth"
)

func TestHealthAndReadiness(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := testRouter(logger, []string{"http://localhost:3000"})
	for _, path := range []string{"/healthz", "/readyz"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Fatalf("%s returned %d", path, recorder.Code)
		}
	}
}

func TestCORSPreflightAllowsConfiguredOrigin(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := testRouter(logger, []string{"https://crm-system-gold-three.vercel.app/"})
	request := httptest.NewRequest(http.MethodOptions, "/api/users/login", nil)
	request.Header.Set("Origin", "https://crm-system-gold-three.vercel.app")
	request.Header.Set("Access-Control-Request-Method", http.MethodPost)
	request.Header.Set("Access-Control-Request-Headers", "content-type,authorization")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("preflight returned %d", recorder.Code)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "https://crm-system-gold-three.vercel.app" {
		t.Fatalf("Access-Control-Allow-Origin = %q", got)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Headers"); got != "content-type,authorization" {
		t.Fatalf("Access-Control-Allow-Headers = %q", got)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Methods"); !strings.Contains(got, http.MethodPost) {
		t.Fatalf("Access-Control-Allow-Methods = %q", got)
	}
}

func TestCORSPreflightRejectsUnknownOrigin(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := testRouter(logger, []string{"https://crm-system-gold-three.vercel.app"})
	request := httptest.NewRequest(http.MethodOptions, "/api/users/login", nil)
	request.Header.Set("Origin", "https://unknown.example")
	request.Header.Set("Access-Control-Request-Method", http.MethodPost)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("preflight returned %d", recorder.Code)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("Access-Control-Allow-Origin = %q", got)
	}
}

func testRouter(logger *slog.Logger, origins []string) http.Handler {
	return Router(
		Modules{}, auth.New("01234567890123456789012345678901", time.Hour),
		nil, logger, origins, func() error { return nil },
	)
}
