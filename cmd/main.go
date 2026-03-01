package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	} else {
		log.Println("Postgres connected!")
	}

	rdb, err := cache.NewRedisClient(config.AppConfig)
	if err != nil {
		log.Fatal("Error connecting to Redis.", err)
	} else {
		log.Println("Redis connected!")
	}

	urlRepo := database.NewUrlDBRepository(postgres)
	urlRedisCache := cache.NewUrlRedisCache(rdb)

	shortenerHandlers := handlers.NewShortenerRoutes(urlRepo, urlRedisCache)

	mux := http.NewServeMux()
	shortenerHandlers.RegisterRoutes(mux)
	h := middlewares.RecoverMiddleware(middlewares.LoggingMiddleware(logger, mux))

	server := http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	go func() {
		log.Println("Server started!")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Shutdown error,", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	fmt.Println("Shutting down...")
	if err = server.Shutdown(ctx); err != nil {
		fmt.Println("Error on shutting down:", err)
	}

	fmt.Println("Server stopped.")
}
