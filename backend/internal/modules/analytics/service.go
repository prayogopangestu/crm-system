package analytics

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/prayogopangestu/crm-system/backend/internal/shared"
)

type Repository interface {
	DashboardStats(ctx context.Context, organizationID string, now time.Time) (DashboardStats, error)
	ConversionChart(ctx context.Context, organizationID string, now time.Time) ([]ConversionPoint, error)
	Activities(ctx context.Context, organizationID string, limit int) ([]Activity, error)
	Leaderboard(ctx context.Context, organizationID string, month time.Time) ([]LeaderboardEntry, error)
	LostReasons(ctx context.Context, organizationID string) ([]LostReason, error)
	Goals(ctx context.Context, organizationID string) ([]PerformanceGoal, error)
}

type Service struct {
	repository Repository
	cache      shared.CacheHelper
	location   *time.Location
}

func NewService(repository Repository, cache shared.CacheHelper, location *time.Location) *Service {
	return &Service{repository: repository, cache: cache, location: location}
}

func (s *Service) DashboardStats(ctx context.Context, principal shared.Principal) (DashboardStats, error) {
	var result DashboardStats
	key := "crm:" + principal.OrganizationID + ":dashboard:stats"
	err := s.cache.Load(ctx, key, time.Minute, &result, func() error {
		var err error
		result, err = s.repository.DashboardStats(ctx, principal.OrganizationID, time.Now().In(s.location))
		return err
	})
	return result, err
}

func (s *Service) ConversionChart(ctx context.Context, principal shared.Principal) ([]ConversionPoint, error) {
	var result []ConversionPoint
	key := "crm:" + principal.OrganizationID + ":dashboard:conversion"
	err := s.cache.Load(ctx, key, time.Minute, &result, func() error {
		var err error
		result, err = s.repository.ConversionChart(ctx, principal.OrganizationID, time.Now().In(s.location))
		return err
	})
	return result, err
}

func (s *Service) Activities(ctx context.Context, principal shared.Principal, limit int) ([]Activity, error) {
	if limit < 1 {
		limit = 5
	}
	if limit > 100 {
		limit = 100
	}
	return s.repository.Activities(ctx, principal.OrganizationID, limit)
}

func (s *Service) Leaderboard(ctx context.Context, principal shared.Principal, period string) ([]LeaderboardEntry, error) {
	month := time.Now().In(s.location)
	if strings.EqualFold(period, "Bulan Lalu") {
		month = month.AddDate(0, -1, 0)
	}
	var result []LeaderboardEntry
	key := fmt.Sprintf("crm:%s:reports:leaderboard:%s", principal.OrganizationID, month.Format("2006-01"))
	err := s.cache.Load(ctx, key, 5*time.Minute, &result, func() error {
		var err error
		result, err = s.repository.Leaderboard(ctx, principal.OrganizationID, month)
		return err
	})
	return result, err
}

func (s *Service) LostReasons(ctx context.Context, principal shared.Principal) ([]LostReason, error) {
	var result []LostReason
	key := "crm:" + principal.OrganizationID + ":reports:lost-reasons"
	err := s.cache.Load(ctx, key, 5*time.Minute, &result, func() error {
		var err error
		result, err = s.repository.LostReasons(ctx, principal.OrganizationID)
		return err
	})
	return result, err
}

func (s *Service) Goals(ctx context.Context, principal shared.Principal) ([]PerformanceGoal, error) {
	var result []PerformanceGoal
	key := "crm:" + principal.OrganizationID + ":reports:goals"
	err := s.cache.Load(ctx, key, 5*time.Minute, &result, func() error {
		var err error
		result, err = s.repository.Goals(ctx, principal.OrganizationID)
		return err
	})
	return result, err
}

func (s *Service) ExportCSV(ctx context.Context, principal shared.Principal) ([]byte, error) {
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

func (s *Service) ExportPDF(ctx context.Context, principal shared.Principal) ([]byte, error) {
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
