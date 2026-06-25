package http

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	nethttp "net/http"
	"strconv"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"github.com/prayogopangestu/crm-system/backend/internal/usecase"
	"github.com/prayogopangestu/crm-system/backend/pkg/response"
)

type Handler struct {
	services Services
	logger   *slog.Logger
}

type Services struct {
	Users         usecase.UserUseCase
	Contacts      usecase.ContactUseCase
	Deals         usecase.DealUseCase
	Tasks         usecase.TaskUseCase
	Analytics     usecase.AnalyticsUseCase
	Settings      usecase.SettingsUseCase
	Search        usecase.SearchUseCase
	Notifications usecase.NotificationUseCase
}

func NewHandler(services Services, logger *slog.Logger) *Handler {
	return &Handler{services: services, logger: logger}
}

func decodeJSON(w nethttp.ResponseWriter, r *nethttp.Request, dst any) bool {
	r.Body = nethttp.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		message := "JSON request tidak valid"
		if errors.Is(err, io.EOF) {
			message = "request body wajib diisi"
		}
		response.Error(w, nethttp.StatusBadRequest, "invalid_json", message, requestIDFromContext(r.Context()), nil)
		return false
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		response.Error(w, nethttp.StatusBadRequest, "invalid_json", "request hanya boleh memiliki satu objek JSON", requestIDFromContext(r.Context()), nil)
		return false
	}
	return true
}

func principal(r *nethttp.Request) domain.Principal {
	value, _ := domain.PrincipalFromContext(r.Context())
	return value
}

func parseInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func writeError(w nethttp.ResponseWriter, r *nethttp.Request, err error) {
	status := nethttp.StatusInternalServerError
	code := "internal_error"
	message := "terjadi kesalahan internal"
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		status, code, message = nethttp.StatusBadRequest, "validation_error", "request tidak valid"
	case errors.Is(err, domain.ErrUnauthorized):
		status, code, message = nethttp.StatusUnauthorized, "unauthorized", "email atau password tidak valid"
	case errors.Is(err, domain.ErrForbidden):
		status, code, message = nethttp.StatusForbidden, "forbidden", "Anda tidak memiliki akses"
	case errors.Is(err, domain.ErrNotFound):
		status, code, message = nethttp.StatusNotFound, "not_found", "data tidak ditemukan"
	case errors.Is(err, domain.ErrConflict):
		status, code, message = nethttp.StatusConflict, "conflict", "data sudah digunakan"
	case errors.Is(err, domain.ErrStageInUse):
		status, code, message = nethttp.StatusConflict, "stage_in_use", "tahapan masih digunakan oleh deal"
	case errors.Is(err, domain.ErrInviteExpired):
		status, code, message = nethttp.StatusGone, "invite_expired", "undangan telah kedaluwarsa"
	case errors.Is(err, domain.ErrInviteUsed):
		status, code, message = nethttp.StatusConflict, "invite_used", "undangan telah digunakan"
	}
	if status == nethttp.StatusInternalServerError {
		slog.Error("request failed", "error", err, "request_id", requestIDFromContext(r.Context()))
	}
	response.Error(w, status, code, message, requestIDFromContext(r.Context()), nil)
}
