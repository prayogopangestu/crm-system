package usecase

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/prayogopangestu/crm-system/backend/internal/domain"
)

func (s *Service) DashboardStats(ctx context.Context, principal domain.Principal) (domain.DashboardStats, error) {
	var result domain.DashboardStats
	key := "crm:" + principal.OrganizationID + ":dashboard:stats"
	err := s.cached(ctx, key, time.Minute, &result, func() error {
		var err error
		result, err = s.analytics.DashboardStats(ctx, principal.OrganizationID, time.Now().In(s.location))
		return err
	})
	return result, err
}

func (s *Service) ConversionChart(ctx context.Context, principal domain.Principal) ([]domain.ConversionPoint, error) {
	var result []domain.ConversionPoint
	key := "crm:" + principal.OrganizationID + ":dashboard:conversion"
	err := s.cached(ctx, key, time.Minute, &result, func() error {
		var err error
		result, err = s.analytics.ConversionChart(ctx, principal.OrganizationID, time.Now().In(s.location))
		return err
	})
	return result, err
}

func (s *Service) Activities(ctx context.Context, principal domain.Principal, limit int) ([]domain.Activity, error) {
	if limit < 1 {
		limit = 5
	}
	if limit > 100 {
		limit = 100
	}
	return s.analytics.Activities(ctx, principal.OrganizationID, limit)
}

func (s *Service) Leaderboard(ctx context.Context, principal domain.Principal, period string) ([]domain.LeaderboardEntry, error) {
	month := time.Now().In(s.location)
	if strings.EqualFold(period, "Bulan Lalu") {
		month = month.AddDate(0, -1, 0)
	}
	var result []domain.LeaderboardEntry
	key := fmt.Sprintf("crm:%s:reports:leaderboard:%s", principal.OrganizationID, month.Format("2006-01"))
	err := s.cached(ctx, key, 5*time.Minute, &result, func() error {
		var err error
		result, err = s.analytics.Leaderboard(ctx, principal.OrganizationID, month)
		return err
	})
	return result, err
}

func (s *Service) LostReasons(ctx context.Context, principal domain.Principal) ([]domain.LostReason, error) {
	var result []domain.LostReason
	key := "crm:" + principal.OrganizationID + ":reports:lost-reasons"
	err := s.cached(ctx, key, 5*time.Minute, &result, func() error {
		var err error
		result, err = s.analytics.LostReasons(ctx, principal.OrganizationID)
		return err
	})
	return result, err
}

func (s *Service) Goals(ctx context.Context, principal domain.Principal) ([]domain.PerformanceGoal, error) {
	var result []domain.PerformanceGoal
	key := "crm:" + principal.OrganizationID + ":reports:goals"
	err := s.cached(ctx, key, 5*time.Minute, &result, func() error {
		var err error
		result, err = s.analytics.Goals(ctx, principal.OrganizationID)
		return err
	})
	return result, err
}

func (s *Service) Search(ctx context.Context, principal domain.Principal, query string) (domain.SearchResult, error) {
	query = strings.TrimSpace(query)
	if len(query) < 2 {
		return domain.SearchResult{}, domain.ErrInvalidInput
	}
	sum := sha256.Sum256([]byte(strings.ToLower(query)))
	key := "crm:" + principal.OrganizationID + ":search:" + hex.EncodeToString(sum[:8])
	var result domain.SearchResult
	err := s.cached(ctx, key, 30*time.Second, &result, func() error {
		var err error
		result, err = s.search.Search(ctx, principal.OrganizationID, query)
		return err
	})
	return result, err
}

func (s *Service) ExportCSV(ctx context.Context, principal domain.Principal) ([]byte, error) {
	leaders, err := s.Leaderboard(ctx, principal, "Bulan Ini")
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	buffer.WriteString("\xEF\xBB\xBF")
	writer := csv.NewWriter(&buffer)
	_ = writer.Write([]string{"Peringkat", "Nama", "Peran", "Pendapatan", "Tren"})
	for _, item := range leaders {
		_ = writer.Write([]string{
			strconv.Itoa(item.Rank), item.Name, item.Role, strconv.FormatInt(item.Amount, 10), item.Trend,
		})
	}
	writer.Flush()
	return buffer.Bytes(), writer.Error()
}

func (s *Service) ExportPDF(ctx context.Context, principal domain.Principal) ([]byte, error) {
	leaders, err := s.Leaderboard(ctx, principal, "Bulan Ini")
	if err != nil {
		return nil, err
	}
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 12, "Laporan Performa Penjualan")
	pdf.Ln(16)
	pdf.SetFont("Arial", "B", 11)
	for _, header := range []string{"Rank", "Nama", "Peran", "Pendapatan", "Tren"} {
		pdf.CellFormat(38, 8, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 10)
	for _, item := range leaders {
		for _, value := range []string{
			strconv.Itoa(item.Rank), item.Name, item.Role, strconv.FormatInt(item.Amount, 10), item.Trend,
		} {
			pdf.CellFormat(38, 8, value, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}
	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
