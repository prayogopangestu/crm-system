package http

import (
	"fmt"
	nethttp "net/http"
	"time"

	"github.com/prayogopangestu/crm-system/backend/pkg/response"
)

func (h *Handler) dashboardStats(w nethttp.ResponseWriter, r *nethttp.Request) {
	result, err := h.services.Analytics.DashboardStats(r.Context(), principal(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, result)
}

func (h *Handler) conversionChart(w nethttp.ResponseWriter, r *nethttp.Request) {
	result, err := h.services.Analytics.ConversionChart(r.Context(), principal(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, result)
}

func (h *Handler) activities(w nethttp.ResponseWriter, r *nethttp.Request) {
	result, err := h.services.Analytics.Activities(r.Context(), principal(r), parseInt(r.URL.Query().Get("limit"), 5))
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, result)
}

func (h *Handler) leaderboard(w nethttp.ResponseWriter, r *nethttp.Request) {
	result, err := h.services.Analytics.Leaderboard(r.Context(), principal(r), r.URL.Query().Get("period"))
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, result)
}

func (h *Handler) lostReasons(w nethttp.ResponseWriter, r *nethttp.Request) {
	result, err := h.services.Analytics.LostReasons(r.Context(), principal(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, result)
}

func (h *Handler) goals(w nethttp.ResponseWriter, r *nethttp.Request) {
	result, err := h.services.Analytics.Goals(r.Context(), principal(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.JSON(w, nethttp.StatusOK, result)
}

func (h *Handler) exportCSV(w nethttp.ResponseWriter, r *nethttp.Request) {
	data, err := h.services.Analytics.ExportCSV(r.Context(), principal(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeDownload(w, "text/csv; charset=utf-8", "laporan-penjualan-"+time.Now().Format("2006-01-02")+".csv", data)
}

func (h *Handler) exportPDF(w nethttp.ResponseWriter, r *nethttp.Request) {
	data, err := h.services.Analytics.ExportPDF(r.Context(), principal(r))
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeDownload(w, "application/pdf", "laporan-penjualan-"+time.Now().Format("2006-01-02")+".pdf", data)
}

func writeDownload(w nethttp.ResponseWriter, contentType, filename string, data []byte) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	w.WriteHeader(nethttp.StatusOK)
	_, _ = w.Write(data)
}
