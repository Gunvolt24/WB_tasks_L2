package repo

import (
	"time"

	"github.com/Gunvolt24/wb_l2/L2.18/internal/domain"
)

type EventRepository interface {
	CreateEvent(event *domain.Event) error
	UpdateEvent(event *domain.Event) error
	DeleteEvent(eventID string) error
	GetEventsForPeriod(userID string, start, end time.Time) ([]*domain.Event, error)
}
