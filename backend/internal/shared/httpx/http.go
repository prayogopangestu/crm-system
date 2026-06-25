package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/prayogopangestu/crm-system/backend/internal/shared"
	"github.com/prayogopangestu/crm-system/backend/pkg/response"
)

type requestIDKey struct{}

func WithRequestID(r *http.Request, value string) *http.Request {
	return r.WithContext(contextWithRequestID(r, value))
}

func contextWithRequestID(r *http.Request, value string) context.Context {
	return context.WithValue(r.Context(), requestIDKey{}, value)
}

func RequestID(ctx context.Context) string {
	value, _ := ctx.Value(requestIDKey{}).(string)
	return value
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		message := "JSON request tidak valid"
		if errors.Is(err, io.EOF) {
			message = "request body wajib diisi"
		}
		response.Error(w, http.StatusBadRequest, "invalid_json", message, RequestID(r.Context()), nil)
		return false
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		response.Error(w, http.StatusBadRequest, "invalid_json", "request hanya boleh memiliki satu objek JSON", RequestID(r.Context()), nil)
		return false
	}
	return true
}

func Principal(r *http.Request) shared.Principal {
	value, _ := shared.PrincipalFromContext(r.Context())
	return value
}

func ParseInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func WriteError(logger *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	status := http.StatusInternalServerError
	code := "internal_error"
	message := "terjadi kesalahan internal"
	switch {
	case errors.Is(err, shared.ErrInvalidInput):
		status, code, message = http.StatusBadRequest, "validation_error", "request tidak valid"
	case errors.Is(err, shared.ErrUnauthorized):
		status, code, message = http.StatusUnauthorized, "unauthorized", "email atau password tidak valid"
	case errors.Is(err, shared.ErrForbidden):
		status, code, message = http.StatusForbidden, "forbidden", "Anda tidak memiliki akses"
	case errors.Is(err, shared.ErrNotFound):
		status, code, message = http.StatusNotFound, "not_found", "data tidak ditemukan"
	case errors.Is(err, shared.ErrConflict):
		status, code, message = http.StatusConflict, "conflict", "data sudah digunakan"
	case errors.Is(err, shared.ErrStageInUse):
		status, code, message = http.StatusConflict, "stage_in_use", "tahapan masih digunakan oleh deal"
	case errors.Is(err, shared.ErrInviteExpired):
		status, code, message = http.StatusGone, "invite_expired", "undangan telah kedaluwarsa"
	case errors.Is(err, shared.ErrInviteUsed):
		status, code, message = http.StatusConflict, "invite_used", "undangan telah digunakan"
	}
	if status == http.StatusInternalServerError {
		logger.Error("request failed", "error", err, "request_id", RequestID(r.Context()))
	}
	response.Error(w, status, code, message, RequestID(r.Context()), nil)
}
