package database

import (
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBConfig interface {
	DBUser() string
	DBPass() string
	DBHost() string
	DBPort() int
	DBName() string
}

func GetPostgresDsn(cfg DBConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser(),
		cfg.DBPass(),
		cfg.DBHost(),
		cfg.DBPort(),
		cfg.DBName(),
	)
}

func ConnectPostgres(cfg DBConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", GetPostgresDsn(cfg))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
