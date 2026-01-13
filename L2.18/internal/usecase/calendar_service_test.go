package usecase

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Gunvolt24/wb_l2/L2.18/internal/domain"
)

func TestCalendarService_AddEvent_Success(t *testing.T) {
	fakeRepo := &fakeEventRepository{
		CreateEventFunc: func(event *domain.Event) error {
			if event.UserID == "" {
				return fmt.Errorf("userID is empty")
			}
			return nil
		},
	}

	service := &CalendarServiceImpl{repo: fakeRepo}

	eventID, err := service.AddEvent(
		"123",
		"Test event",
		"Description",
		time.Now().Add(time.Hour),
		time.Now().Add(2*time.Hour),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if eventID == "" {
		t.Fatal("expected non-empty eventID")
	}
}

func TestCalendarService_AddEvent_Fail(t *testing.T) {
	fakeRepo := &fakeEventRepository{}

	service := &CalendarServiceImpl{repo: fakeRepo}

	now := time.Now()

	tests := []struct {
		name      string
		userID    string
		title     string
		desc      string
		startTime time.Time
		endTime   time.Time
		wantErr   string
	}{
		{"Empty userID", "", "Title", "Desc", now, now.Add(time.Hour), "user_id must be specified"},
		{"End before start", "1", "Title", "Desc", now, now.Add(-time.Hour), "end_time must be after start_time"},
		{"Title empty", "1", "", "Desc", now, now.Add(time.Hour), "title must be specified"},
		{"Title too long", "1", string(make([]byte, MaxTitleLength+1)), "Desc", now, now.Add(time.Hour), "title length must be < 100 characters"},
		{"Description empty", "1", "Title", "", now, now.Add(time.Hour), "description must be specified"},
		{"Description too long", "1", "Title", string(make([]byte, MaxDescriptionLength+1)), now, now.Add(time.Hour), "description length must be < 500 characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.AddEvent(tt.userID, tt.title, tt.desc, tt.startTime, tt.endTime)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestCalendarService_UpdateEvent_Success(t *testing.T) {
	fakeRepo := &fakeEventRepository{
		UpdateEventFunc: func(event *domain.Event) error {
			if event.EventID == "" {
				return fmt.Errorf("eventID is empty")
			}
			return nil
		},
	}

	service := &CalendarServiceImpl{repo: fakeRepo}

	err := service.UpdateEvent(
		"1234",
		"Update Test event",
		"Updated Description",
		time.Now().Add(time.Hour),
		time.Now().Add(3*time.Hour),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCalendarService_UpdateEvent_Fail(t *testing.T) {
	fakeRepo := &fakeEventRepository{}

	service := &CalendarServiceImpl{repo: fakeRepo}

	now := time.Now()

	tests := []struct {
		name      string
		eventID   string
		title     string
		desc      string
		startTime time.Time
		endTime   time.Time
		wantErr   string
	}{
		{"Empty eventID", "", "Title", "Desc", now, now.Add(time.Hour), "event_id must be specified"},
		{"End before start", "1", "Title", "Desc", now, now.Add(-time.Hour), "end_time must be after start_time"},
		{"Title empty", "1", "", "Desc", now, now.Add(time.Hour), "title must be specified"},
		{"Title too long", "1", string(make([]byte, MaxTitleLength+1)), "Desc", now, now.Add(time.Hour), "title length must be < 100 characters"},
		{"Description empty", "1", "Title", "", now, now.Add(time.Hour), "description must be specified"},
		{"Description too long", "1", "Title", string(make([]byte, MaxDescriptionLength+1)), now, now.Add(time.Hour), "description length must be < 500 characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.UpdateEvent(tt.eventID, tt.title, tt.desc, tt.startTime, tt.endTime)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestCalendarService_DeleteEvent_Success(t *testing.T) {
	fakeRepo := &fakeEventRepository{
		DeleteEventFunc: func(eventID string) error {
			if eventID == "" {
				return fmt.Errorf("eventID is empty")
			}
			return nil
		},
	}

	service := &CalendarServiceImpl{repo: fakeRepo}

	err := service.DeleteEvent(
		"1234",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCalendarService_DeleteEvent_Fail(t *testing.T) {
	fakeRepo := &fakeEventRepository{}

	service := &CalendarServiceImpl{repo: fakeRepo}

	err := service.DeleteEvent("")

	if err == nil {
		t.Fatalf("expected error for empty eventID, got nil")
	}
}

func TestCalendarService_GetEventsForDay(t *testing.T) {
	now := time.Now()

	fakeEvents := []*domain.Event{
		{EventID: "1", Title: "Test", StartTime: now, EndTime: now.Add(time.Hour)},
	}

	fakeRepo := &fakeEventRepository{
		GetEventsForPeriodFunc: func(userID string, start, end time.Time) ([]*domain.Event, error) {
			if userID == "fail" {
				return nil, fmt.Errorf("repo error")
			}
			return fakeEvents, nil
		},
	}

	service := &CalendarServiceImpl{repo: fakeRepo}

	tests := []struct {
		name    string
		userID  string
		date    string
		wantErr string
		wantLen int
	}{
		{"Empty userID", "", "2025-12-19", "user_id must be specified", 0},
		{"Bad date", "1", "19-12-2025", "date format must be YYYY-MM-DD", 0},
		{"Repo error", "fail", "2025-12-19", "repo error", 0},
		{"Success", "1", "2025-12-19", "", len(fakeEvents)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := service.GetEventsForDay(tt.userID, tt.date)

			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error %q, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(events) != tt.wantLen {
				t.Fatalf("expected %d events, got %d", tt.wantLen, len(events))
			}
		})
	}
}

func TestCalendarService_GetEventsForWeek(t *testing.T) {
	now := time.Date(2025, 12, 19, 0, 0, 0, 0, time.UTC)
	fakeEvents := []*domain.Event{
		{EventID: "1", Title: "Week Event", StartTime: now, EndTime: now.Add(time.Hour)},
	}

	var capturedStart, capturedEnd time.Time

	fakeRepo := &fakeEventRepository{
		GetEventsForPeriodFunc: func(userID string, start, end time.Time) ([]*domain.Event, error) {
			capturedStart = start
			capturedEnd = end
			if userID == "fail" {
				return nil, fmt.Errorf("repo error")
			}
			return fakeEvents, nil
		},
	}

	service := &CalendarServiceImpl{repo: fakeRepo}

	tests := []struct {
		name      string
		userID    string
		date      string
		wantErr   string
		wantLen   int
		wantStart time.Time
		wantEnd   time.Time
	}{
		{"Empty userID", "", "2025-12-19", "user_id must be specified", 0, time.Time{}, time.Time{}},
		{"Bad date", "1", "19-12-2025", "date format must be YYYY-MM-DD", 0, time.Time{}, time.Time{}},
		{"Repo error", "fail", "2025-12-19", "repo error", 0, time.Time{}, time.Time{}},
		{"Success", "1", "2025-12-19", "", len(fakeEvents),
			time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC),
			time.Date(2025, 12, 22, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := service.GetEventsForWeek(tt.userID, tt.date)

			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error %q, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(events) != tt.wantLen {
				t.Fatalf("expected %d events, got %d", tt.wantLen, len(events))
			}

			if !capturedStart.Equal(tt.wantStart) {
				t.Fatalf("expected start %v, got %v", tt.wantStart, capturedStart)
			}
			if !capturedEnd.Equal(tt.wantEnd) {
				t.Fatalf("expected end %v, got %v", tt.wantEnd, capturedEnd)
			}
		})
	}
}

func TestCalendarService_GetEventsForMonth(t *testing.T) {
	now := time.Date(2025, 12, 19, 0, 0, 0, 0, time.UTC)
	fakeEvents := []*domain.Event{
		{EventID: "1", Title: "Month Event", StartTime: now, EndTime: now.Add(time.Hour)},
	}

	var capturedStart, capturedEnd time.Time

	fakeRepo := &fakeEventRepository{
		GetEventsForPeriodFunc: func(userID string, start, end time.Time) ([]*domain.Event, error) {
			capturedStart = start
			capturedEnd = end
			if userID == "fail" {
				return nil, fmt.Errorf("repo error")
			}
			return fakeEvents, nil
		},
	}

	service := &CalendarServiceImpl{repo: fakeRepo}

	tests := []struct {
		name      string
		userID    string
		date      string
		wantErr   string
		wantLen   int
		wantStart time.Time
		wantEnd   time.Time
	}{
		{"Empty userID", "", "2025-12-19", "user_id must be specified", 0, time.Time{}, time.Time{}},
		{"Bad date", "1", "19-12-2025", "date format must be YYYY-MM-DD", 0, time.Time{}, time.Time{}},
		{"Repo error", "fail", "2025-12-19", "repo error", 0, time.Time{}, time.Time{}},
		{"Success", "1", "2025-12-19", "", len(fakeEvents),
			time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := service.GetEventsForMonth(tt.userID, tt.date)

			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error %q, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(events) != tt.wantLen {
				t.Fatalf("expected %d events, got %d", tt.wantLen, len(events))
			}

			if !capturedStart.Equal(tt.wantStart) {
				t.Fatalf("expected start %v, got %v", tt.wantStart, capturedStart)
			}
			if !capturedEnd.Equal(tt.wantEnd) {
				t.Fatalf("expected end %v, got %v", tt.wantEnd, capturedEnd)
			}
		})
	}
}

type fakeEventRepository struct {
	CreateEventFunc        func(event *domain.Event) error
	UpdateEventFunc        func(event *domain.Event) error
	DeleteEventFunc        func(eventID string) error
	GetEventsForPeriodFunc func(userID string, start, end time.Time) ([]*domain.Event, error)
}

func (f *fakeEventRepository) CreateEvent(event *domain.Event) error {
	if f.CreateEventFunc != nil {
		return f.CreateEventFunc(event)
	}
	return nil
}

func (f *fakeEventRepository) UpdateEvent(event *domain.Event) error {
	if f.UpdateEventFunc != nil {
		return f.UpdateEventFunc(event)
	}
	return nil
}

func (f *fakeEventRepository) DeleteEvent(eventID string) error {
	if f.DeleteEventFunc != nil {
		return f.DeleteEventFunc(eventID)
	}
	return nil
}

func (f *fakeEventRepository) GetEventsForPeriod(userID string, start, end time.Time) ([]*domain.Event, error) {
	if f.GetEventsForPeriodFunc != nil {
		return f.GetEventsForPeriodFunc(userID, start, end)
	}
	return nil, nil
}

func TestCalendarService_AddEvent_RepoError(t *testing.T) {
	fakeRepo := &fakeEventRepository{
		CreateEventFunc: func(event *domain.Event) error {
			return fmt.Errorf("memory error")
		},
	}

	service := &CalendarServiceImpl{repo: fakeRepo}
	_, err := service.AddEvent("123", "Test", "Desc", time.Now().Add(time.Hour), time.Now().Add(2*time.Hour))
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestCalendarService_UpdateEvent_RepoError(t *testing.T) {
	fakeRepo := &fakeEventRepository{
		UpdateEventFunc: func(event *domain.Event) error {
			return fmt.Errorf("memory error")
		},
	}

	service := &CalendarServiceImpl{repo: fakeRepo}
	err := service.UpdateEvent("123", "Test", "Desc", time.Now().Add(time.Hour), time.Now().Add(2*time.Hour))
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestCalendarService_DeleteEvent_RepoError(t *testing.T) {
	fakeRepo := &fakeEventRepository{
		DeleteEventFunc: func(eventID string) error {
			return fmt.Errorf("memory error")
		},
	}

	service := &CalendarServiceImpl{repo: fakeRepo}
	err := service.DeleteEvent("123")
	if err == nil {
		t.Fatal("expected error from repo")
	}
}
