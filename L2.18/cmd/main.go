package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Gunvolt24/wb_l2/L2.18/internal/transport/rest"
	"github.com/Gunvolt24/wb_l2/L2.18/internal/usecase"
	"github.com/Gunvolt24/wb_l2/L2.18/middleware"
	"github.com/Gunvolt24/wb_l2/L2.18/repo"
	"github.com/joho/godotenv"
)

func main() {
	// Погружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using default environment variables")
	}

	// Default HOST и PORT
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Адрес сервера
	addr := fmt.Sprintf("%s:%s", host, port)

	// Создаем слои приложения
	repository := repo.NewInMemoryRepo()
	service := usecase.NewCalendarService(repository)
	handler := rest.NewHandler(service)
	router := rest.NewRouter(handler)
	loggedRouter := middleware.LoggingMiddleware(router)

	// Общий mux для статики и API
	mux := http.NewServeMux()

	// Статика
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	// REST API эндпоинты
	mux.Handle("/create_event", loggedRouter)
	mux.Handle("/events_for_day", loggedRouter)
	mux.Handle("/events_for_week", loggedRouter)
	mux.Handle("/events_for_month", loggedRouter)

	log.Printf("Server started on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
