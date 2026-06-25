package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"github.com/prayogopangestu/crm-system/backend/pkg/response"
)

func (h *Handler) listTeam(w nethttp.ResponseWriter, r *nethttp.Request) {
	result, err := h.service.ListTeam(r.Context(), principal(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, result)
}

func (h *Handler) inviteMember(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request domain.InviteInput
	if !decodeJSON(w, r, &request) {
		return
	}
	result, err := h.service.InviteMember(r.Context(), principal(r), request)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusCreated, map[string]any{
		"success": true, "data": result.User, "inviteUrl": result.InviteURL,
	})
}

func (h *Handler) revokeMember(w nethttp.ResponseWriter, r *nethttp.Request) {
	if err := h.service.RevokeMember(r.Context(), principal(r), chi.URLParam(r, "id")); err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"success": true, "message": "Anggota tim dihapus"})
}

func (h *Handler) listStages(w nethttp.ResponseWriter, r *nethttp.Request) {
	result, err := h.service.ListStages(r.Context(), principal(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, result)
}

func (h *Handler) createStage(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request domain.StageInput
	if !decodeJSON(w, r, &request) {
		return
	}
	result, err := h.service.CreateStage(r.Context(), principal(r), request)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusCreated, map[string]any{"success": true, "data": result})
}

func (h *Handler) reorderStages(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request struct {
		StagesOrder []string `json:"stagesOrder"`
	}
	if !decodeJSON(w, r, &request) {
		return
	}
	if err := h.service.ReorderStages(r.Context(), principal(r), request.StagesOrder); err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"success": true})
}

func (h *Handler) deleteStage(w nethttp.ResponseWriter, r *nethttp.Request) {
	if err := h.service.DeleteStage(r.Context(), principal(r), chi.URLParam(r, "id")); err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"success": true, "message": "Tahapan dihapus"})
}

func (h *Handler) getTelegram(w nethttp.ResponseWriter, r *nethttp.Request) {
	result, err := h.service.GetTelegram(r.Context(), principal(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, result)
}

func (h *Handler) updateTelegram(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request domain.TelegramInput
	if !decodeJSON(w, r, &request) {
		return
	}
	result, err := h.service.UpdateTelegram(r.Context(), principal(r), request)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"success": true, "enabled": result.Enabled})
}

func (h *Handler) testTelegram(w nethttp.ResponseWriter, r *nethttp.Request) {
	if err := h.service.TestTelegram(r.Context(), principal(r)); err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"success": true, "message": "Pesan uji coba terkirim"})
}

func (h *Handler) search(w nethttp.ResponseWriter, r *nethttp.Request) {
	result, err := h.service.Search(r.Context(), principal(r), r.URL.Query().Get("q"))
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, result)
}

func (h *Handler) notifications(w nethttp.ResponseWriter, r *nethttp.Request) {
	result, err := h.service.Notifications(r.Context(), principal(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, result)
}

func (h *Handler) readNotification(w nethttp.ResponseWriter, r *nethttp.Request) {
	if err := h.service.ReadNotification(r.Context(), principal(r), chi.URLParam(r, "id")); err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"success": true})
}

func (h *Handler) readAllNotifications(w nethttp.ResponseWriter, r *nethttp.Request) {
	if err := h.service.ReadAllNotifications(r.Context(), principal(r)); err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"success": true})
}
