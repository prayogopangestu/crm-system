package http

import (
	"log/slog"
	nethttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"github.com/prayogopangestu/crm-system/backend/pkg/auth"
)

func Router(handler *Handler, tokens *auth.Manager, cache domain.Cache, logger *slog.Logger, origins []string, ready func() error) nethttp.Handler {
	router := chi.NewRouter()
	router.Use(requestID)
	router.Use(chimiddleware.RealIP)
	router.Use(recoverer(logger))
	router.Use(accessLog(logger))
	router.Use(cors(origins))

	router.Get("/healthz", handler.health)
	router.Get("/readyz", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		handler.readyCheck(w, r, ready)
	})

	authLimiter := rateLimit(cache, logger, "auth", 5, time.Minute)
	router.With(authLimiter).Post("/api/auth/login", handler.login)
	router.With(authLimiter).Post("/api/auth/register", handler.register)
	router.With(authLimiter).Post("/api/auth/accept-invite", handler.acceptInvite)

	router.Group(func(protected chi.Router) {
		protected.Use(authenticate(tokens))
		protected.Get("/api/profile", handler.getProfile)
		protected.Put("/api/profile", handler.updateProfile)

		protected.Get("/api/dashboard/stats", handler.dashboardStats)
		protected.Get("/api/dashboard/conversion-chart", handler.conversionChart)
		protected.Get("/api/dashboard/activities", handler.activities)

		protected.Get("/api/contacts", handler.listContacts)
		protected.Post("/api/contacts", handler.createContact)
		protected.Put("/api/contacts/{id}", handler.updateContact)
		protected.Delete("/api/contacts/{id}", handler.deleteContact)

		protected.Get("/api/deals", handler.listDeals)
		protected.Post("/api/deals", handler.createDeal)
		protected.Patch("/api/deals/{id}/stage", handler.updateDealStage)
		protected.Put("/api/deals/{id}", handler.updateDeal)
		protected.Delete("/api/deals/{id}", handler.deleteDeal)

		protected.Get("/api/tasks", handler.listTasks)
		protected.Post("/api/tasks", handler.createTask)
		protected.Put("/api/tasks/{id}", handler.updateTask)
		protected.Patch("/api/tasks/{id}/toggle", handler.toggleTask)
		protected.Delete("/api/tasks/{id}", handler.deleteTask)

		protected.Get("/api/reports/leaderboard", handler.leaderboard)
		protected.Get("/api/reports/lost-reasons", handler.lostReasons)
		protected.Get("/api/reports/goals", handler.goals)
		protected.Get("/api/reports/export/csv", handler.exportCSV)
		protected.Get("/api/reports/export/pdf", handler.exportPDF)

		protected.Get("/api/team", handler.listTeam)
		protected.Post("/api/team/invite", handler.inviteMember)
		protected.Delete("/api/team/{id}", handler.revokeMember)
		protected.Get("/api/pipeline-stages", handler.listStages)
		protected.Post("/api/pipeline-stages", handler.createStage)
		protected.Put("/api/pipeline-stages/reorder", handler.reorderStages)
		protected.Delete("/api/pipeline-stages/{id}", handler.deleteStage)

		protected.Get("/api/integrations/telegram", handler.getTelegram)
		protected.Put("/api/integrations/telegram", handler.updateTelegram)
		protected.Post("/api/integrations/telegram/test", handler.testTelegram)

		protected.Get("/api/search", handler.search)
		protected.Get("/api/notifications", handler.notifications)
		protected.Patch("/api/notifications/{id}/read", handler.readNotification)
		protected.Patch("/api/notifications/read-all", handler.readAllNotifications)
	})
	return router
}
