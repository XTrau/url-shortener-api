package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"urlshortener/internal/database"
	"urlshortener/internal/handlers"
	"urlshortener/internal/middlewares"
)

func main() {
	textHandler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(textHandler)

	db, err := database.NewPostresDB()
	if err != nil {
		log.Fatal("Error on creating Postres connection pool.", err)
	}

	urlRepo := database.NewUrlDBRepository(db)
	shortenerHandlers := handlers.NewShortenerRoutes(urlRepo)

	mux := http.NewServeMux()
	shortenerHandlers.RegisterRoutes(mux)
	server := middlewares.RecoverMiddleware(middlewares.LoggingMiddleware(logger, mux))

	logger.Info("Server started!")
	http.ListenAndServe(":8080", server)
}
