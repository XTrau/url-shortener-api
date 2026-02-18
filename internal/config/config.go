package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
)

type Config struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort int
	DBName string
}

var AppConfig Config

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found, using environment variables")
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("Error parsing DB_PORT:", err)
	}

	AppConfig = Config{
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: port,
		DBName: os.Getenv("DB_NAME"),
	}
}

func GetPostgresDsn() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		AppConfig.DBUser,
		AppConfig.DBPass,
		AppConfig.DBHost,
		AppConfig.DBPort,
		AppConfig.DBName,
	)
}
