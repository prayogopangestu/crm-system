package postgres

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
)

func (r *Repository) DashboardStats(ctx context.Context, organizationID string, now time.Time) (domain.DashboardStats, error) {
	var stats domain.DashboardStats
	var currentLeads, previousLeads, currentWon, previousWon, revenue int64
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	previousMonth := monthStart.AddDate(0, -1, 0)
	weekStart := dayStart(now).AddDate(0, 0, -7)
	previousWeek := weekStart.AddDate(0, 0, -7)
	err := r.pool.QueryRow(ctx, `
		SELECT
		  count(*) FILTER (WHERE c.created_at >= $2),
		  count(*) FILTER (WHERE c.created_at >= $3 AND c.created_at < $2),
		  (SELECT count(*) FROM deals d WHERE d.organization_id=$1 AND d.deleted_at IS NULL
		    AND d.stage_key='won' AND d.updated_at >= $4),
		  (SELECT count(*) FROM deals d WHERE d.organization_id=$1 AND d.deleted_at IS NULL
		    AND d.stage_key='won' AND d.updated_at >= $5 AND d.updated_at < $4),
		  (SELECT COALESCE(sum(value),0) FROM deals d WHERE d.organization_id=$1
		    AND d.deleted_at IS NULL AND d.stage_key='won'),
		  (SELECT count(*) FROM tasks t WHERE t.organization_id=$1 AND t.deleted_at IS NULL
		    AND t.completed=false AND t.due_date <= $6::date AND t.priority='Tinggi')
		FROM contacts c WHERE c.organization_id=$1 AND c.deleted_at IS NULL`,
		organizationID, monthStart, previousMonth, weekStart, previousWeek, now,
	).Scan(
		&currentLeads, &previousLeads, &currentWon, &previousWon, &revenue, &stats.UrgentTasksCount,
	)
	if err != nil {
		return stats, err
	}
	if err := r.pool.QueryRow(ctx,
		`SELECT count(*) FROM contacts WHERE organization_id=$1 AND deleted_at IS NULL`,
		organizationID,
	).Scan(&stats.TotalLeads); err != nil {
		return stats, err
	}
	if err := r.pool.QueryRow(ctx,
		`SELECT count(*) FROM deals WHERE organization_id=$1 AND deleted_at IS NULL AND stage_key='won'`,
		organizationID,
	).Scan(&stats.DealWonCount); err != nil {
		return stats, err
	}
	stats.LeadsTrend = trendPercent(currentLeads, previousLeads) + " bulan ini"
	stats.WonTrend = fmt.Sprintf("%+d dari minggu lalu", currentWon-previousWon)
	stats.TotalRevenue = formatRupiah(revenue)
	stats.RevenueTrend = "Stabil"
	return stats, nil
}

func (r *Repository) ConversionChart(ctx context.Context, organizationID string, now time.Time) ([]domain.ConversionPoint, error) {
	start := time.Date(now.Year(), now.Month()-5, 1, 0, 0, 0, 0, now.Location())
	rows, err := r.pool.Query(ctx, `
		WITH months AS (
		  SELECT generate_series($2::date, date_trunc('month',$3::date), interval '1 month') AS month
		)
		SELECT to_char(m.month,'Mon'),
		       CASE WHEN count(d.id)=0 THEN 0
		            ELSE round(100.0*count(d.id) FILTER (WHERE d.stage_key='won')/count(d.id),2)
		       END
		FROM months m
		LEFT JOIN deals d ON d.organization_id=$1 AND d.deleted_at IS NULL
		  AND date_trunc('month',d.created_at)=m.month
		GROUP BY m.month ORDER BY m.month`,
		organizationID, start, now,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	points := make([]domain.ConversionPoint, 0, 6)
	for rows.Next() {
		var point domain.ConversionPoint
		if err := rows.Scan(&point.Name, &point.Conversion); err != nil {
			return nil, err
		}
		points = append(points, point)
	}
	return points, rows.Err()
}

func (r *Repository) Activities(ctx context.Context, organizationID string, limit int) ([]domain.Activity, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id,actor_name,action,target,is_highlight,created_at
		FROM activities WHERE organization_id=$1
		ORDER BY created_at DESC LIMIT $2`,
		organizationID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	now := time.Now().In(r.location)
	items := make([]domain.Activity, 0)
	for rows.Next() {
		var item domain.Activity
		if err := rows.Scan(
			&item.ID, &item.User, &item.Action, &item.Target, &item.IsHighlight, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		item.Time = humanTime(item.CreatedAt, now)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) Leaderboard(ctx context.Context, organizationID string, month time.Time) ([]domain.LeaderboardEntry, error) {
	start := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	end := start.AddDate(0, 1, 0)
	previous := start.AddDate(0, -1, 0)
	rows, err := r.pool.Query(ctx, `
		SELECT trim(u.first_name || ' ' || u.last_name),u.role,u.avatar_url,
		       COALESCE(sum(d.value) FILTER (WHERE d.updated_at >= $2 AND d.updated_at < $3),0),
		       COALESCE(sum(d.value) FILTER (WHERE d.updated_at >= $4 AND d.updated_at < $2),0)
		FROM users u
		LEFT JOIN deals d ON d.assignee_id=u.id AND d.organization_id=u.organization_id
		  AND d.deleted_at IS NULL AND d.stage_key='won'
		WHERE u.organization_id=$1 AND u.revoked_at IS NULL
		GROUP BY u.id ORDER BY 4 DESC`,
		organizationID, start, end, previous,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.LeaderboardEntry, 0)
	for rows.Next() {
		var item domain.LeaderboardEntry
		var previousAmount int64
		if err := rows.Scan(&item.Name, &item.Role, &item.AvatarURL, &item.Amount, &previousAmount); err != nil {
			return nil, err
		}
		item.Rank = len(items) + 1
		trend := percentChange(item.Amount, previousAmount)
		item.IsPositive = trend >= 0
		item.Trend = fmt.Sprintf("%+.0f%%", trend)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) LostReasons(ctx context.Context, organizationID string) ([]domain.LostReason, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT COALESCE(NULLIF(lost_reason,''),'Tidak diketahui'),count(*)
		FROM deals
		WHERE organization_id=$1 AND deleted_at IS NULL AND stage_key='lost'
		GROUP BY 1 ORDER BY 2 DESC`,
		organizationID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.LostReason, 0)
	var total int64
	for rows.Next() {
		var item domain.LostReason
		if err := rows.Scan(&item.Name, &item.Value); err != nil {
			return nil, err
		}
		total += item.Value
		items = append(items, item)
	}
	colors := []string{"#6366f1", "#f59e0b", "#10b981", "#ef4444", "#8b5cf6"}
	for index := range items {
		if total > 0 {
			items[index].Percentage = int64(math.Round(100 * float64(items[index].Value) / float64(total)))
		}
		items[index].Color = colors[index%len(colors)]
	}
	return items, rows.Err()
}

func (r *Repository) Goals(ctx context.Context, organizationID string) ([]domain.PerformanceGoal, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT g.month,g.goal,
		       COALESCE(sum(d.value) FILTER (
		         WHERE d.updated_at >= g.month
		           AND d.updated_at < g.month + interval '1 month'
		       ),0)
		FROM performance_goals g
		LEFT JOIN deals d ON d.organization_id=g.organization_id
		  AND d.deleted_at IS NULL AND d.stage_key='won'
		WHERE g.organization_id=$1
		GROUP BY g.id ORDER BY g.month DESC`,
		organizationID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.PerformanceGoal, 0)
	for rows.Next() {
		var month time.Time
		var item domain.PerformanceGoal
		if err := rows.Scan(&month, &item.Goal, &item.Actual); err != nil {
			return nil, err
		}
		item.Month = indonesianMonth(month.Month())
		if item.Goal > 0 {
			item.Percentage = int64(math.Round(100 * float64(item.Actual) / float64(item.Goal)))
		}
		if item.Percentage >= 100 {
			item.Status = fmt.Sprintf("Tercapai (%d%%)", item.Percentage)
		} else {
			item.Status = fmt.Sprintf("Kurang (%d%%)", item.Percentage)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func trendPercent(current, previous int64) string {
	return fmt.Sprintf("%+.0f%%", percentChange(current, previous))
}

func percentChange(current, previous int64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return 100 * float64(current-previous) / float64(previous)
}

func indonesianMonth(month time.Month) string {
	names := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	return names[int(month)-1]
}
