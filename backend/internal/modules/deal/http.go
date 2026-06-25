package deal

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
	router.Get("/api/deals", h.list)
	router.Post("/api/deals", h.create)
	router.Patch("/api/deals/{id}/stage", h.updateStage)
	router.Put("/api/deals/{id}", h.update)
	router.Delete("/api/deals/{id}", h.delete)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.List(r.Context(), httpx.Principal(r))
	if err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var request Input
	if !httpx.DecodeJSON(w, r, &request) {
		return
	}
	result, err := h.service.Create(r.Context(), httpx.Principal(r), request)
	if err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{"success": true, "data": result})
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var request Input
	if !httpx.DecodeJSON(w, r, &request) {
		return
	}
	result, err := h.service.Update(r.Context(), httpx.Principal(r), chi.URLParam(r, "id"), request)
	if err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"success": true, "data": result})
}

func (h *Handler) updateStage(w http.ResponseWriter, r *http.Request) {
	var request StageInput
	if !httpx.DecodeJSON(w, r, &request) {
		return
	}
	if err := h.service.UpdateStage(r.Context(), httpx.Principal(r), chi.URLParam(r, "id"), request); err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"success": true, "message": "Tahap deal berhasil diperbarui"})
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.Context(), httpx.Principal(r), chi.URLParam(r, "id")); err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"success": true, "message": "Deal berhasil dihapus"})
}
