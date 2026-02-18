package database

import (
	"database/sql"
	"urlshortener/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostresDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", config.GetPostgresDsn())
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
