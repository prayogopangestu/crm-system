package search

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
	router.Get("/api/search", h.search)
}

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.Search(r.Context(), httpx.Principal(r), r.URL.Query().Get("q"))
	if err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}
