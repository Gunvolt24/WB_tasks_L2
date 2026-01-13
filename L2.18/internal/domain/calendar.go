package domain

import "time"

type Event struct {
	UserID      string    `json:"user_id"`
	EventID     string    `json:"event_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
}
