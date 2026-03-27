package config

import (
	"log/slog"

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
