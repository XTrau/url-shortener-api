package config

import (
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
)

type Config struct {
	LogLevel slog.Level

	DBUser string
	DBPass string
	DBHost string
	DBPort int
	DBName string

	RedisHost     string
	RedisPort     int
	RedisUser     string
	RedisPassword string
	RedisDatabase int
}

var AppConfig Config

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found, using environment variables")
	}

	var logLevel slog.Level
	logLevelStr := os.Getenv("LOG_LEVEL")

	if strings.EqualFold("debug", logLevelStr) {
		logLevel = slog.LevelDebug
	} else if strings.EqualFold("error", logLevelStr) {
		logLevel = slog.LevelError
	} else if strings.EqualFold("error", logLevelStr) {
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

	AppConfig = Config{
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
