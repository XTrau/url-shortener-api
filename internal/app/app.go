package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"urlshortener/internal/cache"
	"urlshortener/internal/config"
	"urlshortener/internal/database"
	"urlshortener/internal/handlers"
	"urlshortener/internal/middlewares"

	"github.com/joho/godotenv"
)

var AppConfig config.Config

func init() {
	if err := godotenv.Load(); err == nil {
		slog.Info(".env file loaded.")
	}

	var logLevel slog.Level
	logLevelStr := os.Getenv("LOG_LEVEL")

	if strings.EqualFold("debug", logLevelStr) {
		logLevel = slog.LevelDebug
	} else if strings.EqualFold("error", logLevelStr) {
		logLevel = slog.LevelError
	} else if strings.EqualFold("warning", logLevelStr) {
		logLevel = slog.LevelWarn
	} else {
		logLevel = slog.LevelInfo
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("Error parsing DB_PORT:", err)
	}

	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Fatal("Error parsing REDIS_PORT:", err)
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatal("Error parsing REDIS_DB:", err)
	}

	AppConfig = config.Config{
		LogLevel: logLevel,

		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: dbPort,
		DBName: os.Getenv("DB_NAME"),

		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     redisPort,
		RedisUser:     os.Getenv("REDIS_USER"),
		RedisPassword: os.Getenv("REDIS_PASS"),
		RedisDatabase: redisDB,
	}
}

func Run() error {
	logOpts := &slog.HandlerOptions{
		Level: AppConfig.LogLevel,
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, logOpts)))

	postgres, err := database.NewPostresDB(AppConfig)
	if err != nil {
		return fmt.Errorf("Error on creating Postres connection pool: %w", err)
	}

	slog.Info("Postgres connected!")

	err = database.RunMigrations(AppConfig)
	if err != nil {
		return fmt.Errorf("Error on running Postgres migrations: %w", err)
	}

	rdb, err := database.NewRedisClient(AppConfig)
	if err != nil {
		return fmt.Errorf("Error connecting to Redis: %w", err)
	}

	slog.Info("Redis connected!")

	rkvs := cache.NewRedisKeyValueStorage(rdb, time.Second)

	urlRepo := database.NewUrlDBRepository(postgres)
	urlCache := cache.NewUrlCache(rkvs)

	mux := http.NewServeMux()
	r := handlers.NewShortenerRoutes(urlRepo, urlCache)
	r.RegisterRoutes(mux)

	wh := handlers.NewWebHandlers()
	wh.RegisterRoutes(mux)

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
