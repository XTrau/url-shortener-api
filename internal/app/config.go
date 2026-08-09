package app

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	logLevel slog.Level

	dbUser string
	dbPass string
	dbHost string
	dbPort int
	dbName string

	redisHost     string
	redisPort     string
	redisUser     string
	redisPassword string
	redisDatabase int
}

func (c *Config) LogLevel() slog.Level { return c.logLevel }

func (c *Config) DBUser() string { return c.dbUser }
func (c *Config) DBPass() string { return c.dbPass }
func (c *Config) DBHost() string { return c.dbHost }
func (c *Config) DBPort() int    { return c.dbPort }
func (c *Config) DBName() string { return c.dbName }

func (c *Config) RedisHost() string     { return c.redisHost }
func (c *Config) RedisPort() string     { return c.redisPort }
func (c *Config) RedisUser() string     { return c.redisUser }
func (c *Config) RedisPassword() string { return c.redisPassword }
func (c *Config) RedisDatabase() int    { return c.redisDatabase }

func LoadConfig() (*Config, error) {
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
		return nil, fmt.Errorf("Error parsing DB_PORT: %v", err)
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, fmt.Errorf("Error parsing REDIS_DB: %v", err)
	}

	return &Config{
		logLevel: logLevel,

		dbUser: os.Getenv("DB_USER"),
		dbPass: os.Getenv("DB_PASS"),
		dbHost: os.Getenv("DB_HOST"),
		dbPort: dbPort,
		dbName: os.Getenv("DB_NAME"),

		redisHost:     os.Getenv("REDIS_HOST"),
		redisPort:     os.Getenv("REDIS_PORT"),
		redisUser:     os.Getenv("REDIS_USER"),
		redisPassword: os.Getenv("REDIS_PASS"),
		redisDatabase: redisDB,
	}, nil
}
