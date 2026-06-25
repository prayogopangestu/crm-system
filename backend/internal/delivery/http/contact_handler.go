package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"github.com/prayogopangestu/crm-system/backend/pkg/response"
)

func (h *Handler) listContacts(w nethttp.ResponseWriter, r *nethttp.Request) {
	result, err := h.services.Contacts.ListContacts(
		r.Context(), principal(r), r.URL.Query().Get("search"), r.URL.Query().Get("status"),
		parseInt(r.URL.Query().Get("page"), 1), parseInt(r.URL.Query().Get("limit"), 20),
	)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, result)
}

func (h *Handler) createContact(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request domain.ContactInput
	if !decodeJSON(w, r, &request) {
		return
	}
	result, err := h.services.Contacts.CreateContact(r.Context(), principal(r), request)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusCreated, map[string]any{"success": true, "data": result})
}

func (h *Handler) updateContact(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request domain.ContactInput
	if !decodeJSON(w, r, &request) {
		return
	}
	result, err := h.services.Contacts.UpdateContact(r.Context(), principal(r), chi.URLParam(r, "id"), request)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"success": true, "data": result})
}

func (h *Handler) deleteContact(w nethttp.ResponseWriter, r *nethttp.Request) {
	if err := h.services.Contacts.DeleteContact(r.Context(), principal(r), chi.URLParam(r, "id")); err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"success": true, "message": "Kontak berhasil dihapus"})
}
