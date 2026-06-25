package http

import (
	nethttp "net/http"

	"github.com/prayogopangestu/crm-system/backend/pkg/response"
)

func (h *Handler) health(w nethttp.ResponseWriter, _ *nethttp.Request) {
	response.JSON(w, nethttp.StatusOK, map[string]any{"status": "ok"})
}

func (h *Handler) readyCheck(w nethttp.ResponseWriter, r *nethttp.Request, check func() error) {
	if err := check(); err != nil {
		response.JSON(w, nethttp.StatusServiceUnavailable, map[string]any{"status": "not_ready"})
		return
	}
	response.JSON(w, nethttp.StatusOK, map[string]any{"status": "ready"})
}
