package httpserver

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/analytics"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/contact"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/deal"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/integration"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/notification"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/pipeline"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/search"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/task"
	"github.com/prayogopangestu/crm-system/backend/internal/modules/user"
	"github.com/prayogopangestu/crm-system/backend/internal/shared"
	"github.com/prayogopangestu/crm-system/backend/internal/shared/httpx"
	"github.com/prayogopangestu/crm-system/backend/pkg/auth"
	"github.com/prayogopangestu/crm-system/backend/pkg/response"
)

type Modules struct {
	User         *user.Handler
	Contact      *contact.Handler
	Deal         *deal.Handler
	Task         *task.Handler
	Analytics    *analytics.Handler
	Pipeline     *pipeline.Handler
	Integration  *integration.Handler
	Notification *notification.Handler
	Search       *search.Handler
}

func Router(modules Modules, tokens *auth.Manager, cache shared.Cache, logger *slog.Logger, origins []string, ready func() error) http.Handler {
	router := chi.NewRouter()
	router.Use(requestID)
	router.Use(chimiddleware.RealIP)
	router.Use(recoverer(logger))
	router.Use(accessLog(logger))
	router.Use(cors(origins))

	router.Get("/healthz", health)
	router.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := ready(); err != nil {
			logger.Warn("ready probe dependency check failed", "error", err)
		}
		response.JSON(w, http.StatusOK, map[string]any{"status": "ready"})
	})

	authLimiter := rateLimit(cache, logger, "auth", 5, time.Minute)
	if modules.User != nil {
		modules.User.PublicRoutes(router, authLimiter)
	}

	router.Group(func(protected chi.Router) {
		protected.Use(authenticate(tokens))
		if modules.User != nil {
			modules.User.ProtectedRoutes(protected)
		}
		if modules.Contact != nil {
			modules.Contact.Routes(protected)
		}
		if modules.Deal != nil {
			modules.Deal.Routes(protected)
		}
		if modules.Task != nil {
			modules.Task.Routes(protected)
		}
		if modules.Analytics != nil {
			modules.Analytics.Routes(protected)
		}
		if modules.Pipeline != nil {
			modules.Pipeline.Routes(protected)
		}
		if modules.Integration != nil {
			modules.Integration.Routes(protected)
		}
		if modules.Notification != nil {
			modules.Notification.Routes(protected)
		}
		if modules.Search != nil {
			modules.Search.Routes(protected)
		}
	})
	return router
}

func health(w http.ResponseWriter, _ *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := r.Header.Get("X-Request-ID")
		if value == "" {
			value = uuid.NewString()
		}
		w.Header().Set("X-Request-ID", value)
		next.ServeHTTP(w, httpx.WithRequestID(r, value))
	})
}

func recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if value := recover(); value != nil {
					logger.Error("panic recovered", "panic", value, "request_id", httpx.RequestID(r.Context()))
					response.Error(w, http.StatusInternalServerError, "internal_error", "terjadi kesalahan internal", httpx.RequestID(r.Context()), nil)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func accessLog(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("http request",
				"method", r.Method, "path", r.URL.Path,
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", httpx.RequestID(r.Context()),
			)
		})
	}
}

func cors(origins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(origins))
	for _, origin := range origins {
		allowed[strings.TrimSpace(origin)] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if allowed[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func authenticate(tokens *auth.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				response.Error(w, http.StatusUnauthorized, "unauthorized", "Bearer token diperlukan", httpx.RequestID(r.Context()), nil)
				return
			}
			principal, err := tokens.Parse(strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")))
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "unauthorized", "token tidak valid atau kedaluwarsa", httpx.RequestID(r.Context()), nil)
				return
			}
			next.ServeHTTP(w, r.WithContext(shared.WithPrincipal(r.Context(), principal)))
		})
	}
}

func rateLimit(cache shared.Cache, logger *slog.Logger, route string, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cache == nil {
				next.ServeHTTP(w, r)
				return
			}
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				host = r.RemoteAddr
			}
			bucket := time.Now().Unix() / int64(window.Seconds())
			key := "crm:rate:" + route + ":" + host + ":" + time.Unix(bucket*int64(window.Seconds()), 0).Format("20060102150405")
			allowed, err := cache.Allow(r.Context(), key, limit, window+time.Minute)
			if err != nil {
				logger.Warn("rate limiter unavailable", "error", err)
				next.ServeHTTP(w, r)
				return
			}
			if !allowed {
				response.Error(w, http.StatusTooManyRequests, "rate_limited", "terlalu banyak percobaan, coba lagi nanti", httpx.RequestID(r.Context()), nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
