package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"github.com/prayogopangestu/crm-system/backend/pkg/response"
)

func (h *Handler) listDeals(w nethttp.ResponseWriter, r *nethttp.Request) {
	result, err := h.service.ListDeals(r.Context(), principal(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, result)
}

func (h *Handler) createDeal(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request domain.DealInput
	if !decodeJSON(w, r, &request) {
		return
	}
	result, err := h.service.CreateDeal(r.Context(), principal(r), request)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusCreated, map[string]any{"success": true, "data": result})
}

func (h *Handler) updateDeal(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request domain.DealInput
	if !decodeJSON(w, r, &request) {
		return
	}
	result, err := h.service.UpdateDeal(r.Context(), principal(r), chi.URLParam(r, "id"), request)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"success": true, "data": result})
}

func (h *Handler) updateDealStage(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request domain.StageUpdateInput
	if !decodeJSON(w, r, &request) {
		return
	}
	if err := h.service.UpdateDealStage(r.Context(), principal(r), chi.URLParam(r, "id"), request); err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"success": true, "message": "Tahap deal berhasil diperbarui"})
}

func (h *Handler) deleteDeal(w nethttp.ResponseWriter, r *nethttp.Request) {
	if err := h.service.DeleteDeal(r.Context(), principal(r), chi.URLParam(r, "id")); err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"success": true, "message": "Deal berhasil dihapus"})
}
