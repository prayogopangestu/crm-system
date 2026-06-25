package http

import (
	nethttp "net/http"
	"strings"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"github.com/prayogopangestu/crm-system/backend/pkg/response"
)

func (h *Handler) register(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request struct {
		Name        string `json:"name"`
		FullName    string `json:"fullName"`
		CompanyName string `json:"companyName"`
		Email       string `json:"email"`
		Password    string `json:"password"`
	}
	if !decodeJSON(w, r, &request) {
		return
	}
	name := strings.TrimSpace(request.Name)
	if name == "" {
		name = strings.TrimSpace(request.FullName)
	}
	_, err := h.services.Users.Register(r.Context(), domain.RegisterInput{
		Name: name, CompanyName: request.CompanyName, Email: request.Email, Password: request.Password,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusCreated, map[string]any{"success": true, "message": "Registrasi berhasil"})
}

func (h *Handler) login(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request domain.LoginInput
	if !decodeJSON(w, r, &request) {
		return
	}
	result, err := h.services.Users.Login(r.Context(), request)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, result)
}

func (h *Handler) acceptInvite(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if !decodeJSON(w, r, &request) {
		return
	}
	_, err := h.services.Users.AcceptInvite(r.Context(), request.Token, request.Password)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"success": true, "message": "Undangan berhasil diaktifkan"})
}

func (h *Handler) getProfile(w nethttp.ResponseWriter, r *nethttp.Request) {
	result, err := h.services.Users.Profile(r.Context(), principal(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, result)
}

func (h *Handler) updateProfile(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request domain.UpdateProfileInput
	if !decodeJSON(w, r, &request) {
		return
	}
	result, err := h.services.Users.UpdateProfile(r.Context(), principal(r), request)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{
		"success": true, "message": "Profil diperbarui", "data": result,
	})
}
