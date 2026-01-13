package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Gunvolt24/wb_l2/L2.18/internal/domain"
	"github.com/Gunvolt24/wb_l2/L2.18/internal/dto"
)

// DTO для запросов Add, Update, Delete, GetEvents
type AddEventRequest struct {
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

type UpdateEventRequest struct {
	EventID     string `json:"event_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

type DeleteEventRequest struct {
	EventID string `json:"event_id"`
}

type GetEventsRequest struct {
	UserID string
	Date   string
}

// ParseAddEventRequest парсит HTTP запрос для добавления события
func ParseAddEventRequest(r *http.Request) (*dto.AddEventInput, error) {
	var req AddEventRequest
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Println("failed to close request body:", err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid json: %w", err)
	}

	if req.UserID == "" {
		return nil, domain.NewValidationError("user_id is required")
	}

	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, domain.NewValidationError("start_time must be RFC3339")
	}

	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, domain.NewValidationError("end_time must be RFC3339")
	}

	return &dto.AddEventInput{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   start,
		EndTime:     end,
	}, nil
}

// ParseUpdateEventRequest парсит HTTP запрос для обновления события
func ParseUpdateEventRequest(r *http.Request) (*dto.UpdateEventInput, error) {
	var req UpdateEventRequest
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Println("failed to close request body:", err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid json: %w", err)
	}

	if req.EventID == "" {
		return nil, domain.NewValidationError("event_id is required")
	}

	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, domain.NewValidationError("start_time must be RFC3339")
	}

	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, domain.NewValidationError("end_time must be RFC3339")
	}

	return &dto.UpdateEventInput{
		EventID:     req.EventID,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   start,
		EndTime:     end,
	}, nil
}

// ParseDeleteEventRequest парсит HTTP запрос для удаления события
func ParseDeleteEventRequest(r *http.Request) (*dto.DeleteEventInput, error) {
	var req DeleteEventRequest
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Println("failed to close request body:", err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid json: %w", err)
	}

	if req.EventID == "" {
		return nil, domain.NewValidationError("event_id is required")
	}

	return &dto.DeleteEventInput{
		EventID: req.EventID,
	}, nil
}

// ParseGetEventsForDayRequest парсит HTTP запрос для получения событий за день
func ParseGetEventsForDayRequest(r *http.Request) (*dto.GetEventsInput, error) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		return nil, domain.NewValidationError("user_id is required")
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		return nil, domain.NewValidationError("date is required")
	}

	const layout = "2006-01-02"
	_, err := time.Parse(layout, dateStr)
	if err != nil {
		return nil, domain.NewValidationError("date must be in format YYYY-MM-DD")
	}

	return &dto.GetEventsInput{
		UserID: userID,
		Date:   dateStr,
	}, nil
}

// Для простоты парсинга week/month используем ту же логику, что и для day
func ParseGetEventsForWeekRequest(r *http.Request) (*dto.GetEventsInput, error) {
	return ParseGetEventsForDayRequest(r)
}

func ParseGetEventsForMonthRequest(r *http.Request) (*dto.GetEventsInput, error) {
	return ParseGetEventsForDayRequest(r)
}
