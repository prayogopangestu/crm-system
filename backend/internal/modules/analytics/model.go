package analytics

import "time"

type DashboardStats struct {
	TotalLeads       int64  `json:"totalLeads"`
	LeadsTrend       string `json:"leadsTrend"`
	DealWonCount     int64  `json:"dealWonCount"`
	WonTrend         string `json:"wonTrend"`
	TotalRevenue     string `json:"totalRevenue"`
	RevenueTrend     string `json:"revenueTrend"`
	UrgentTasksCount int64  `json:"urgentTasksCount"`
}

type ConversionPoint struct {
	Name       string  `json:"name"`
	Conversion float64 `json:"Konversi"`
}

type Activity struct {
	ID          string    `json:"id"`
	User        string    `json:"user"`
	Action      string    `json:"action"`
	Target      string    `json:"target"`
	Time        string    `json:"time"`
	IsHighlight bool      `json:"isHighlight"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}

type LeaderboardEntry struct {
	Rank       int    `json:"rank"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	Amount     int64  `json:"amount"`
	Trend      string `json:"trend"`
	IsPositive bool   `json:"isPositive"`
	AvatarURL  string `json:"avatarUrl"`
}

type LostReason struct {
	Name       string `json:"name"`
	Value      int64  `json:"value"`
	Percentage int64  `json:"percentage"`
	Color      string `json:"color"`
}

type PerformanceGoal struct {
	Month      string `json:"month"`
	Goal       int64  `json:"goal"`
	Actual     int64  `json:"actual"`
	Status     string `json:"status"`
	Percentage int64  `json:"percentage"`
}
