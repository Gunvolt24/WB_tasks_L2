package usecase

import (
	"time"

	"github.com/Gunvolt24/wb_l2/L2.18/internal/domain"
	"github.com/Gunvolt24/wb_l2/L2.18/repo"
	"github.com/google/uuid"
)

const MaxTitleLength = 100
const MaxDescriptionLength = 500

type CalendarService interface {
	AddEvent(userID, title, description string, startTime, endTime time.Time) (string, error)
	UpdateEvent(eventID, title, description string, startTime, endTime time.Time) error
	DeleteEvent(eventID string) error
	GetEventsForDay(userID string, date string) ([]*domain.Event, error)
	GetEventsForWeek(userID string, date string) ([]*domain.Event, error)
	GetEventsForMonth(userID string, date string) ([]*domain.Event, error)
}

type CalendarServiceImpl struct {
	repo repo.EventRepository
}

func NewCalendarService(repo repo.EventRepository) CalendarService {
	return &CalendarServiceImpl{repo: repo}
}

func (s *CalendarServiceImpl) AddEvent(userID, title, description string, startTime, endTime time.Time) (string, error) {
	// Проверям, что userID не пустой
	if userID == "" {
		return "", domain.NewValidationError("user_id must be specified")
	}

	// Проверка времени
	if !endTime.After(startTime) {
		return "", domain.NewBusinessError("end_time must be after start_time")
	}

	// Проверка, что заголовок не пустой и не превышает максимального количества символов
	if title == "" {
		return "", domain.NewValidationError("title must be specified")
	}

	if len(title) > MaxTitleLength {
		return "", domain.NewValidationError("title length must be < 100 characters")
	}

	// Проверка, что описание не пустое и не превышает максимального количества символов
	if description == "" {
		return "", domain.NewValidationError("description must be specified")
	}

	if len(description) > MaxDescriptionLength {
		return "", domain.NewValidationError("description length must be < 500 characters")
	}

	// Генерируем eventUUID и CreatedAt
	eventUUID := uuid.New()
	now := time.Now()

	// Присваиваем значения
	event := &domain.Event{
		UserID:      userID,
		Title:       title,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
		EventID:     eventUUID.String(),
		CreatedAt:   now,
	}

	// Передаем событие в репозиторий
	if err := s.repo.CreateEvent(event); err != nil {
		return "", domain.NewBusinessError(err.Error())
	}

	return event.EventID, nil

}

func (s *CalendarServiceImpl) UpdateEvent(eventID, title, description string, startTime, endTime time.Time) error {
	// Проверям, что eventID не пустой
	if eventID == "" {
		return domain.NewValidationError("event_id must be specified")
	}

	// Проверка времени
	if !endTime.After(startTime) {
		return domain.NewBusinessError("end_time must be after start_time")
	}

	// Проверка, что заголовок не пустой и не превышает максимального количества символов
	if title == "" {
		return domain.NewValidationError("title must be specified")
	}

	if len(title) > MaxTitleLength {
		return domain.NewValidationError("title length must be < 100 characters")
	}

	// Проверка, что описание не пустое и не превышает максимального количества символов
	if description == "" {
		return domain.NewValidationError("description must be specified")
	}

	if len(description) > MaxDescriptionLength {
		return domain.NewValidationError("description length must be < 500 characters")
	}

	// Присваиваем значения
	event := &domain.Event{
		EventID:     eventID,
		Title:       title,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	// Передаем событие в репозиторий
	if err := s.repo.UpdateEvent(event); err != nil {
		return domain.NewBusinessError(err.Error())
	}

	return nil
}

func (s *CalendarServiceImpl) DeleteEvent(eventID string) error {
	// Проверям, что eventID не пустой
	if eventID == "" {
		return domain.NewValidationError("event_id must be specified")
	}

	// Передаем событие в репозиторий
	if err := s.repo.DeleteEvent(eventID); err != nil {
		return domain.NewBusinessError(err.Error())
	}

	return nil
}

func (s *CalendarServiceImpl) GetEventsForDay(userID string, date string) ([]*domain.Event, error) {
	// Проверям, что userID не пустой
	if userID == "" {
		return nil, domain.NewValidationError("user_id must be specified")
	}

	// Парсим дату
	const layout = "2006-01-02"
	day, err := time.Parse(layout, date)
	if err != nil {
		return nil, domain.NewValidationError("date format must be YYYY-MM-DD")
	}

	// Вычисляем начало и конец дня
	start := time.Date(
		day.Year(), day.Month(), day.Day(),
		0, 0, 0, 0,
		day.Location(),
	)

	end := start.Add(24 * time.Hour)

	events, err := s.repo.GetEventsForPeriod(userID, start, end)
	if err != nil {
		return nil, domain.NewBusinessError(err.Error())
	}

	return events, nil
}

func (s *CalendarServiceImpl) GetEventsForWeek(userID string, date string) ([]*domain.Event, error) {
	// Проверям, что userID не пустой
	if userID == "" {
		return nil, domain.NewValidationError("user_id must be specified")
	}

	// Парсим дату
	const layout = "2006-01-02"
	day, err := time.Parse(layout, date)
	if err != nil {
		return nil, domain.NewValidationError("date format must be YYYY-MM-DD")
	}

	// Вычисляем начало и конец недели (понедельник - воскресенье)
	weekday := int(day.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	start := time.Date(day.Year(), day.Month(), day.Day()-weekday+1, 0, 0, 0, 0, day.Location())

	end := start.AddDate(0, 0, 7)

	events, err := s.repo.GetEventsForPeriod(userID, start, end)
	if err != nil {
		return nil, domain.NewBusinessError(err.Error())
	}

	return events, nil
}

func (s *CalendarServiceImpl) GetEventsForMonth(userID string, date string) ([]*domain.Event, error) {
	// Проверям, что userID не пустой
	if userID == "" {
		return nil, domain.NewValidationError("user_id must be specified")
	}

	// Парсим дату
	const layout = "2006-01-02"
	day, err := time.Parse(layout, date)
	if err != nil {
		return nil, domain.NewValidationError("date format must be YYYY-MM-DD")
	}

	// Вычисляем начало и конец месяца
	start := time.Date(day.Year(), day.Month(), 1, 0, 0, 0, 0, day.Location())

	end := start.AddDate(0, 1, 0)

	events, err := s.repo.GetEventsForPeriod(userID, start, end)
	if err != nil {
		return nil, domain.NewBusinessError(err.Error())
	}

	return events, nil
}
