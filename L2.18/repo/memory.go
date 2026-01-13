package repo

import (
	"fmt"
	"sync"
	"time"

	"github.com/Gunvolt24/wb_l2/L2.18/internal/domain"
)

// InMemoryEventRepository реализует репозиторий событий в памяти
type InMemoryEventRepository struct {
	mu     sync.RWMutex
	events map[string]*domain.Event
}

// NewInMemoryRepo создает новый экземпляр InMemoryEventRepository
func NewInMemoryRepo() *InMemoryEventRepository {
	return &InMemoryEventRepository{
		events: make(map[string]*domain.Event),
	}
}

// CreateEvent добавляет новое событие в репозиторий
func (r *InMemoryEventRepository) CreateEvent(event *domain.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[event.EventID]; exists {
		return fmt.Errorf("event with ID %s already exists", event.EventID)
	}

	r.events[event.EventID] = event
	return nil
}

// UpdateEvent обновляет существующее событие в репозитории
func (r *InMemoryEventRepository) UpdateEvent(event *domain.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.events[event.EventID]
	if !exists {
		return domain.NewNotFoundError("event is not found")
	}

	// Обновляем только изменяемые поля
	existing.Title = event.Title
	existing.Description = event.Description
	existing.StartTime = event.StartTime
	existing.EndTime = event.EndTime

	return nil
}

// DeleteEvent удаляет событие из репозитория по его ID
func (r *InMemoryEventRepository) DeleteEvent(eventID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[eventID]; !exists {
		return domain.NewNotFoundError("event is not found")
	}

	delete(r.events, eventID)
	return nil
}

// GetEventsForPeriod возвращает события для указанного пользователя в заданном периоде времени
func (r *InMemoryEventRepository) GetEventsForPeriod(userID string, start, end time.Time) ([]*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Event
	for _, event := range r.events {
		if event.UserID == userID && !event.EndTime.Before(start) && !event.StartTime.After(end) {
			result = append(result, event)
		}
	}

	return result, nil
}
