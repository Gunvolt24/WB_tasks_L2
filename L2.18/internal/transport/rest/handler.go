package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Gunvolt24/wb_l2/L2.18/internal/domain"
	"github.com/Gunvolt24/wb_l2/L2.18/internal/dto"
	"github.com/Gunvolt24/wb_l2/L2.18/internal/usecase"
)

type Handler struct {
	service usecase.CalendarService
}

func NewHandler(service usecase.CalendarService) *Handler {
	return &Handler{service: service}
}

// Ответы API
func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, ErrorResponse{Error: err.Error()})
}

// Вспомогательная функция для записи JSON ответа
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("failed to write JSON response:", err)
	}
}

// POST /events
func (h *Handler) handleAddEvent(w http.ResponseWriter, r *http.Request) {
	// Парсим запрос в структуру DTO для бизнес логики
	input, err := ParseAddEventRequest(r)
	if err != nil {
		mapErrorToStatusCode(w, err)
		return
	}

	// Вызываем бизнес логику для создания события. Передаем все необходимые поля
	eventID, err := h.service.AddEvent(input.UserID, input.Title, input.Description, input.StartTime, input.EndTime)
	if err != nil {
		mapErrorToStatusCode(w, err)
		return
	}

	// Если событие успешно создано, возвращаем клиенту ответ с кодом 201
	writeJSON(w, http.StatusOK, SuccessResponse{Result: eventID})
}

// PUT /events
func (h *Handler) handleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	// Парсим запрос в структуру DTO для бизнес логики
	input, err := ParseUpdateEventRequest(r)
	if err != nil {
		mapErrorToStatusCode(w, err)
		return
	}

	// Вызываем бизнес логику для создания события. Передаем все необходимые поля
	err = h.service.UpdateEvent(input.EventID, input.Title, input.Description, input.StartTime, input.EndTime)
	if err != nil {
		mapErrorToStatusCode(w, err)
		return
	}

	// Если событие успешно обновлено, возвращаем клиенту ответ с кодом 200
	writeJSON(w, http.StatusOK, SuccessResponse{Result: "updated"})
}

// DELETE /events?id=...
func (h *Handler) handleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	// Парсим запрос в структуру DTO для бизнес логики
	input, err := ParseDeleteEventRequest(r)
	if err != nil {
		mapErrorToStatusCode(w, err)
		return
	}

	// Вызываем бизнес логику для создания события. Передаем все необходимые поля
	err = h.service.DeleteEvent(input.EventID)
	if err != nil {
		mapErrorToStatusCode(w, err)
		return
	}

	// Если событие успешно удалено, возвращаем клиенту ответ с кодом 200
	writeJSON(w, http.StatusOK, SuccessResponse{Result: "deleted"})
}

func (h *Handler) handleGetEvents(w http.ResponseWriter, r *http.Request, period string) {
	var input *dto.GetEventsInput
	var err error

	// Выбираем какой парсер использовать
	switch period {
	case "day":
		input, err = ParseGetEventsForDayRequest(r)
	case "week":
		input, err = ParseGetEventsForWeekRequest(r)
	case "month":
		input, err = ParseGetEventsForMonthRequest(r)
	default:
		mapErrorToStatusCode(w, domain.NewValidationError(fmt.Sprintf("invalid period: %s", period)))
		return
	}

	if err != nil {
		mapErrorToStatusCode(w, err)
		return
	}

	// Вызываем соответствующий метод usecase
	var events []*domain.Event
	switch period {
	case "day":
		events, err = h.service.GetEventsForDay(input.UserID, input.Date)
	case "week":
		events, err = h.service.GetEventsForWeek(input.UserID, input.Date)
	case "month":
		events, err = h.service.GetEventsForMonth(input.UserID, input.Date)
	}

	if err != nil {
		mapErrorToStatusCode(w, err)
		return
	}

	// Формируем HTTP ответ
	var resp EventsResponse
	for _, event := range events {
		resp.Events = append(resp.Events, EventResponse{
			EventID:     event.EventID,
			UserID:      event.UserID,
			Title:       event.Title,
			Description: event.Description,
			StartTime:   event.StartTime.Format(time.RFC3339),
			EndTime:     event.EndTime.Format(time.RFC3339),
		})
	}

	writeJSON(w, http.StatusOK, resp)
}
