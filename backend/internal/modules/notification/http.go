package notification

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prayogopangestu/crm-system/backend/internal/shared/httpx"
	"github.com/prayogopangestu/crm-system/backend/pkg/response"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Routes(router chi.Router) {
	router.Get("/api/notifications", h.list)
	router.Patch("/api/notifications/{id}/read", h.read)
	router.Patch("/api/notifications/read-all", h.readAll)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.List(r.Context(), httpx.Principal(r))
	if err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) read(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Read(r.Context(), httpx.Principal(r), chi.URLParam(r, "id")); err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"success": true})
}

func (h *Handler) readAll(w http.ResponseWriter, r *http.Request) {
	if err := h.service.ReadAll(r.Context(), httpx.Principal(r)); err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"success": true})
}
