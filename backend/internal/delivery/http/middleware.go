package http

import (
	"context"
	"log/slog"
	"net"
	nethttp "net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"github.com/prayogopangestu/crm-system/backend/pkg/auth"
	"github.com/prayogopangestu/crm-system/backend/pkg/response"
)

type requestIDKey struct{}

func requestIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(requestIDKey{}).(string)
	return value
}

func requestID(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		value := r.Header.Get("X-Request-ID")
		if value == "" {
			value = uuid.NewString()
		}
		w.Header().Set("X-Request-ID", value)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey{}, value)))
	})
}

func recoverer(logger *slog.Logger) func(nethttp.Handler) nethttp.Handler {
	return func(next nethttp.Handler) nethttp.Handler {
		return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			defer func() {
				if value := recover(); value != nil {
					logger.Error("panic recovered", "panic", value, "request_id", requestIDFromContext(r.Context()))
					response.Error(w, nethttp.StatusInternalServerError, "internal_error", "terjadi kesalahan internal", requestIDFromContext(r.Context()), nil)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func accessLog(logger *slog.Logger) func(nethttp.Handler) nethttp.Handler {
	return func(next nethttp.Handler) nethttp.Handler {
		return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("http request",
				"method", r.Method, "path", r.URL.Path,
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", requestIDFromContext(r.Context()),
			)
		})
	}
}

func cors(origins []string) func(nethttp.Handler) nethttp.Handler {
	allowed := make(map[string]bool, len(origins))
	for _, origin := range origins {
		allowed[strings.TrimSpace(origin)] = true
	}
	return func(next nethttp.Handler) nethttp.Handler {
		return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			origin := r.Header.Get("Origin")
			if allowed[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			}
			if r.Method == nethttp.MethodOptions {
				w.WriteHeader(nethttp.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func authenticate(tokens *auth.Manager) func(nethttp.Handler) nethttp.Handler {
	return func(next nethttp.Handler) nethttp.Handler {
		return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				response.Error(w, nethttp.StatusUnauthorized, "unauthorized", "Bearer token diperlukan", requestIDFromContext(r.Context()), nil)
				return
			}
			principal, err := tokens.Parse(strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")))
			if err != nil {
				response.Error(w, nethttp.StatusUnauthorized, "unauthorized", "token tidak valid atau kedaluwarsa", requestIDFromContext(r.Context()), nil)
				return
			}
			next.ServeHTTP(w, r.WithContext(domain.WithPrincipal(r.Context(), principal)))
		})
	}
}

func rateLimit(cache domain.Cache, logger *slog.Logger, route string, limit int, window time.Duration) func(nethttp.Handler) nethttp.Handler {
	return func(next nethttp.Handler) nethttp.Handler {
		return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
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
				response.Error(w, nethttp.StatusTooManyRequests, "rate_limited", "terlalu banyak percobaan, coba lagi nanti", requestIDFromContext(r.Context()), nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
