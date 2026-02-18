package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"urlshortener/internal/cache"
	"urlshortener/internal/config"
	"urlshortener/internal/database"
	"urlshortener/internal/handlers"
	"urlshortener/internal/middlewares"
)

func main() {
	textHandler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(textHandler)

	postgres, err := database.NewPostresDB(config.AppConfig)
	if err != nil {
		log.Fatal("Error on creating Postres connection pool.", err)
	}

	rdb, err := cache.NewRedisClient(config.AppConfig)
	if err != nil {
		log.Fatal("Error connecting to Redis.", err)
	}

	urlRepo := database.NewUrlDBRepository(postgres)
	urlRedisCache := cache.NewUrlRedisCache(rdb)

	shortenerHandlers := handlers.NewShortenerRoutes(urlRepo, urlRedisCache)

	mux := http.NewServeMux()
	shortenerHandlers.RegisterRoutes(mux)
	server := middlewares.RecoverMiddleware(middlewares.LoggingMiddleware(logger, mux))

	logger.Info("Server started!")
	http.ListenAndServe(":8080", server)
}
