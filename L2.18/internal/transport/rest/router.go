package rest

import (
	"fmt"
	"net/http"
)

// NewRouter настраивает маршруты для REST API и возвращает http.Handler
func NewRouter(handler *Handler) http.Handler {
	mux := http.NewServeMux()

	// POST /create_event
	mux.HandleFunc("/create_event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}
		handler.handleAddEvent(w, r)
	})

	// POST /update_event
	mux.HandleFunc("/update_event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}
		handler.handleUpdateEvent(w, r)
	})

	// POST /delete_event
	mux.HandleFunc("/delete_event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}
		handler.handleDeleteEvent(w, r)
	})

	// GET /events_for_day
	mux.HandleFunc("/events_for_day", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}
		handler.handleGetEvents(w, r, "day")
	})

	// GET /events_for_week
	mux.HandleFunc("/events_for_week", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}
		handler.handleGetEvents(w, r, "week")
	})

	// GET /events_for_month
	mux.HandleFunc("/events_for_month", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
			return
		}
		handler.handleGetEvents(w, r, "month")
	})

	return mux
}
