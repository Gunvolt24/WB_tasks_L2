package rest

import (
	"net/http"

	"github.com/Gunvolt24/wb_l2/L2.18/internal/domain"
)

// mapErrorToStatusCode присваивает коды статусов HTTP на основе типа ошибки
func mapErrorToStatusCode(w http.ResponseWriter, err error) {
	if errWithCode, ok := err.(domain.ErrorWithCode); ok {
		writeError(w, errWithCode.Code(), errWithCode)
		return
	}
	writeError(w, http.StatusInternalServerError, err)
}
