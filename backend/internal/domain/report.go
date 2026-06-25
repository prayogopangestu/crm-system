package domain

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
