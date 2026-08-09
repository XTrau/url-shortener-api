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
	"urlshortener/internal/database"
	"urlshortener/internal/handlers"
	"urlshortener/internal/middlewares"
)

func Run() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	// Logger
	logOpts := &slog.HandlerOptions{Level: cfg.LogLevel()}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, logOpts)))

	// DB Connection
	postgres, err := database.ConnectPostgres(cfg)
	if err != nil {
		return fmt.Errorf("Error on creating Postres connection: %w", err)
	}

	// Migrations
	if err := RunMigrations(postgres); err != nil {
		return fmt.Errorf("Error on running migrations: %w", err)
	}

	// Redis connection
	rdb, err := database.NewRedisClient(cfg)
	if err != nil {
		return fmt.Errorf("Error connecting to Redis: %w", err)
	}

	// Cache and repositories
	redisCache := cache.NewRedisCache(rdb, time.Second)

	urlCache := cache.NewUrlCache(redisCache)
	urlRepository := database.NewUrlPostgresRepository(postgres)

	// Handlers
	r := handlers.NewShortenerRoutes(urlRepository, urlCache)
	wh := handlers.NewWebHandlers()

	// Register handlers
	mux := http.NewServeMux()

	r.RegisterRoutes(mux)
	wh.RegisterRoutes(mux)

	h := middlewares.LoggingMiddleware(mux)

	server := http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	// Starting server
	go func() {
		slog.Info("Server started!")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(fmt.Sprintf("Server error: %v", err))
		}
	}()

	// Graceful shutdown
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
