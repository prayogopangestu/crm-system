package user

import (
	"log/slog"
	"net/http"
	"strings"

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

func (h *Handler) PublicRoutes(router chi.Router, middleware ...func(http.Handler) http.Handler) {
	router.With(middleware...).Post("/api/auth/login", h.login)
	router.With(middleware...).Post("/api/auth/register", h.register)
	router.With(middleware...).Post("/api/auth/accept-invite", h.acceptInvite)
}

func (h *Handler) ProtectedRoutes(router chi.Router) {
	router.Get("/api/profile", h.profile)
	router.Put("/api/profile", h.updateProfile)
	router.Get("/api/team", h.listTeam)
	router.Post("/api/team/invite", h.inviteMember)
	router.Delete("/api/team/{id}", h.revokeMember)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name        string `json:"name"`
		FullName    string `json:"fullName"`
		CompanyName string `json:"companyName"`
		Email       string `json:"email"`
		Password    string `json:"password"`
	}
	if !httpx.DecodeJSON(w, r, &request) {
		return
	}
	name := strings.TrimSpace(request.Name)
	if name == "" {
		name = strings.TrimSpace(request.FullName)
	}
	if _, err := h.service.Register(r.Context(), RegisterInput{
		Name: name, CompanyName: request.CompanyName, Email: request.Email, Password: request.Password,
	}); err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{"success": true, "message": "Registrasi berhasil"})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var request LoginInput
	if !httpx.DecodeJSON(w, r, &request) {
		return
	}
	result, err := h.service.Login(r.Context(), request)
	if err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) acceptInvite(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if !httpx.DecodeJSON(w, r, &request) {
		return
	}
	if _, err := h.service.AcceptInvite(r.Context(), request.Token, request.Password); err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"success": true, "message": "Undangan berhasil diaktifkan"})
}

func (h *Handler) profile(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.Profile(r.Context(), httpx.Principal(r))
	if err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	var request UpdateProfileInput
	if !httpx.DecodeJSON(w, r, &request) {
		return
	}
	result, err := h.service.UpdateProfile(r.Context(), httpx.Principal(r), request)
	if err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"success": true, "message": "Profil diperbarui", "data": result})
}

func (h *Handler) listTeam(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.ListTeam(r.Context(), httpx.Principal(r))
	if err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) inviteMember(w http.ResponseWriter, r *http.Request) {
	var request InviteInput
	if !httpx.DecodeJSON(w, r, &request) {
		return
	}
	result, err := h.service.InviteMember(r.Context(), httpx.Principal(r), request)
	if err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{"success": true, "data": result.User, "inviteUrl": result.InviteURL})
}

func (h *Handler) revokeMember(w http.ResponseWriter, r *http.Request) {
	if err := h.service.RevokeMember(r.Context(), httpx.Principal(r), chi.URLParam(r, "id")); err != nil {
		httpx.WriteError(h.logger, w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"success": true, "message": "Anggota tim dihapus"})
}
