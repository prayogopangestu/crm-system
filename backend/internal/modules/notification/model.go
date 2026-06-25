package notification

import "time"

type Notification struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Time      string    `json:"time"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}
