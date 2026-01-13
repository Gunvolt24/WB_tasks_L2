package repo

import (
	"testing"
	"time"

	"github.com/Gunvolt24/wb_l2/L2.18/internal/domain"
)

func TestCreateEvent(t *testing.T) {
	repo := NewInMemoryRepo()
	event := &domain.Event{
		EventID:     "1",
		UserID:      "user1",
		Title:       "Test Event",
		Description: "Test Description",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
	}

	err := repo.CreateEvent(event)
	if err != nil {
		t.Errorf("CreateEvent failed: %v", err)
	}

	err = repo.CreateEvent(event)
	if err == nil {
		t.Error("CreateEvent should fail for duplicate ID")
	}
}

func TestUpdateEvent(t *testing.T) {
	repo := NewInMemoryRepo()
	event := &domain.Event{
		EventID:   "1",
		UserID:    "user1",
		Title:     "Original",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(1 * time.Hour),
	}

	if err := repo.CreateEvent(event); err != nil {
		t.Fatalf("failed to create event: %v", err)
	}
	event.Title = "Updated"
	err := repo.UpdateEvent(event)
	if err != nil {
		t.Errorf("UpdateEvent failed: %v", err)
	}

	err = repo.UpdateEvent(&domain.Event{EventID: "999"})
	if err == nil {
		t.Error("UpdateEvent should fail for non-existent event")
	}
}

func TestDeleteEvent(t *testing.T) {
	repo := NewInMemoryRepo()
	event := &domain.Event{EventID: "1", UserID: "user1"}

	if err := repo.CreateEvent(event); err != nil {
		t.Fatalf("failed to create event: %v", err)
	}
	err := repo.DeleteEvent("1")
	if err != nil {
		t.Errorf("DeleteEvent failed: %v", err)
	}

	err = repo.DeleteEvent("999")
	if err == nil {
		t.Error("DeleteEvent should fail for non-existent event")
	}
}

func TestGetEventsForPeriod(t *testing.T) {
	repo := NewInMemoryRepo()
	start := time.Now()
	mid := start.Add(30 * time.Minute)
	end := start.Add(2 * time.Hour)

	event1 := &domain.Event{
		EventID:   "1",
		UserID:    "user1",
		StartTime: start,
		EndTime:   mid,
	}
	event2 := &domain.Event{
		EventID:   "2",
		UserID:    "user1",
		StartTime: mid,
		EndTime:   end,
	}

	if err := repo.CreateEvent(event1); err != nil {
		t.Fatalf("failed to create event: %v", err)
	}
	if err := repo.CreateEvent(event2); err != nil {
		t.Fatalf("failed to create event: %v", err)
	}

	events, err := repo.GetEventsForPeriod("user1", start, end)
	if err != nil || len(events) != 2 {
		t.Errorf("GetEventsForPeriod failed: got %d events, want 2", len(events))
	}

	events, _ = repo.GetEventsForPeriod("user2", start, end)
	if len(events) != 0 {
		t.Errorf("GetEventsForPeriod should return 0 events for different user")
	}
}
