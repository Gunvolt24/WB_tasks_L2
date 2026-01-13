package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// LoggingMiddleware логирует каждый входящий HTTP запрос с указанием метода и пути
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s %s\n", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
