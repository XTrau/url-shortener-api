package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"urlshortener/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func GetPostgresDsn(cfg config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
}

func NewPostresDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", GetPostgresDsn(cfg))
	if err != nil {
		return nil, err
	}

	slog.Debug("Ping to PostgreSQL")
	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.Exec(`CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		url TEXT UNIQUE NOT NULL,
		slug TEXT UNIQUE NOT NULL
	);`)

	return db, nil
}
