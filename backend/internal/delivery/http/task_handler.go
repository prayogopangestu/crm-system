package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"github.com/prayogopangestu/crm-system/backend/pkg/response"
)

func (h *Handler) listTasks(w nethttp.ResponseWriter, r *nethttp.Request) {
	result, err := h.service.ListTasks(r.Context(), principal(r), r.URL.Query().Get("date"), r.URL.Query().Get("status"))
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, result)
}

func (h *Handler) createTask(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request domain.TaskInput
	if !decodeJSON(w, r, &request) {
		return
	}
	result, err := h.service.CreateTask(r.Context(), principal(r), request)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusCreated, map[string]any{"success": true, "data": result})
}

func (h *Handler) updateTask(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request domain.TaskInput
	if !decodeJSON(w, r, &request) {
		return
	}
	result, err := h.service.UpdateTask(r.Context(), principal(r), chi.URLParam(r, "id"), request)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"success": true, "data": result})
}

func (h *Handler) toggleTask(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request struct {
		Completed bool `json:"completed"`
	}
	if !decodeJSON(w, r, &request) {
		return
	}
	if err := h.service.ToggleTask(r.Context(), principal(r), chi.URLParam(r, "id"), request.Completed); err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"success": true, "completed": request.Completed})
}

func (h *Handler) deleteTask(w nethttp.ResponseWriter, r *nethttp.Request) {
	if err := h.service.DeleteTask(r.Context(), principal(r), chi.URLParam(r, "id")); err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"success": true, "message": "Tugas dihapus"})
}
