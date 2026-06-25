package httpserver

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prayogopangestu/crm-system/backend/pkg/auth"
)

func TestHealthAndReadiness(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := Router(
		Modules{}, auth.New("01234567890123456789012345678901", time.Hour),
		nil, logger, []string{"http://localhost:3000"}, func() error { return nil },
	)
	for _, path := range []string{"/healthz", "/readyz"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Fatalf("%s returned %d", path, recorder.Code)
		}
	}
}
