package integration

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
	router.Get("/api/integrations/telegram", h.get)
	router.Put("/api/integrations/telegram", h.update)
	router.Post("/api/integrations/telegram/test", h.test)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.Get(r.Context(), httpx.Principal(r))
	if err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var request TelegramInput
	if !httpx.DecodeJSON(w, r, &request) {
		return
	}
	result, err := h.service.Update(r.Context(), httpx.Principal(r), request)
	if err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"success": true, "enabled": result.Enabled})
}

func (h *Handler) test(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Test(r.Context(), httpx.Principal(r)); err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"success": true, "message": "Pesan uji coba terkirim"})
}
