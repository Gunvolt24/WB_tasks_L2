package dto

import "time"

type AddEventInput struct {
	UserID      string
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}

type UpdateEventInput struct {
	EventID     string
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}

type DeleteEventInput struct {
	EventID string
}

type GetEventsInput struct {
	UserID string
	Date   string
}
