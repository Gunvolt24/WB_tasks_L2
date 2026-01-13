package rest

// DTO для ответов API
type SuccessResponse struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type EventResponse struct {
	EventID     string `json:"event_id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

type EventsResponse struct {
	Events []EventResponse `json:"events"`
}
