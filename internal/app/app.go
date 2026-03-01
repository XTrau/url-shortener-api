package app

import (
	"context"
	"fmt"
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

func Run() error {
	logOpts := &slog.HandlerOptions{
		Level: config.AppConfig.LogLevel,
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, logOpts)))

	postgres, err := database.NewPostresDB(config.AppConfig)
	if err != nil {
		return fmt.Errorf("Error on creating Postres connection pool: %w", err)
	}

	slog.Info("Postgres connected!")

	err = database.RunMigrations(config.AppConfig)
	if err != nil {
		return fmt.Errorf("Error on running Postgres migrations: %w", err)
	}

	rdb, err := cache.NewRedisClient(config.AppConfig)
	if err != nil {
		return fmt.Errorf("Error connecting to Redis: %w", err)
	}

	slog.Info("Redis connected!")

	mux := http.NewServeMux()

	urlRepo := database.NewUrlDBRepository(postgres)
	urlCache := cache.NewUrlRedisCache(rdb)
	r := handlers.NewShortenerRoutes(urlRepo, urlCache)
	r.RegisterRoutes(mux)

	h := middlewares.LoggingMiddleware(mux)

	server := http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	go func() {
		slog.Info("Server started!")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(fmt.Sprintf("Server error: %v", err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	slog.Info("Shutting down...")
	if err = server.Shutdown(ctx); err != nil {
		return fmt.Errorf("Error on shutting down: %w", err)
	}

	slog.Info("Server stopped.")
	return nil
}
