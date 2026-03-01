package database

import (
	"database/sql"
	"errors"
	"log/slog"
	"urlshortener/internal/apperrors"
)

type UrlRepository interface {
	Create(url string, slug string) error
	GetUrlBySlug(slug string) (string, error)
	GetSlugByUrl(url string) (string, error)
}

type UrlDBRepository struct {
	db *sql.DB
}

func NewUrlDBRepository(db *sql.DB) UrlDBRepository {
	return UrlDBRepository{db}
}

func (repo UrlDBRepository) Create(url string, slug string) error {
	slog.Debug("Inserting url to database", slog.String("url", url), slog.String("slug", slug))

	query := "INSERT INTO urls (url, slug) VALUES ($1, $2)"
	_, err := repo.db.Exec(query, url, slug)

	return err
}

func (repo UrlDBRepository) GetUrlBySlug(slug string) (string, error) {
	slog.Debug("Getting url from database", slog.String("slug", slug))

	query := "SELECT url FROM urls WHERE slug=$1"
	row := repo.db.QueryRow(query, slug)

	var url string
	err := row.Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return "", apperrors.ErrUrlNotFound
	}
	return url, err
}

func (repo UrlDBRepository) GetSlugByUrl(url string) (string, error) {
	slog.Debug("Getting url from database", slog.String("url", url))

	query := "SELECT slug FROM urls WHERE url=$1"
	row := repo.db.QueryRow(query, url)

	var slug string
	err := row.Scan(&slug)
	if errors.Is(err, sql.ErrNoRows) {
		return "", apperrors.ErrSlugNotFound
	}
	return slug, err
}
