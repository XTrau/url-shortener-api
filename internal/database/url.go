package database

import (
	"database/sql"
	"errors"
	"log/slog"
	"urlshortener/internal/apperrors"
)

type UrlRepository interface {
	Create(url string, slug string, slugID int) error
	GetUrlBySlug(slug string, slugID int) (string, error)
	GetSlugByUrl(url string) (string, int, error)
	GetFreeSlugID(slug string) (int, error)
}

type UrlDBRepository struct {
	db *sql.DB
}

func NewUrlDBRepository(db *sql.DB) UrlDBRepository {
	return UrlDBRepository{db}
}

func (repo UrlDBRepository) Create(url string, slug string, slugID int) error {
	slog.Debug("Inserting url to database", slog.String("url", url), slog.String("slug", slug), slog.Int("slugID", slugID))

	query := "INSERT INTO urls (url, slug, slug_id) VALUES ($1, $2, $3)"
	_, err := repo.db.Exec(query, url, slug, slugID)

	return err
}

func (repo UrlDBRepository) GetUrlBySlug(slug string, slugID int) (string, error) {
	slog.Debug("Getting url from database", slog.String("slug", slug), slog.Int("slugID", slugID))

	query := "SELECT url FROM urls WHERE slug=$1 AND slug_id=$2"
	row := repo.db.QueryRow(query, slug, slugID)

	var url string
	err := row.Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return "", apperrors.ErrUrlNotFound
	}
	return url, err
}

func (repo UrlDBRepository) GetSlugByUrl(url string) (string, int, error) {
	slog.Debug("Getting url from database", slog.String("url", url))

	query := "SELECT slug, slug_id FROM urls WHERE url=$1"
	row := repo.db.QueryRow(query, url)

	var slug string
	var slugID int

	err := row.Scan(&slug, &slugID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", 0, apperrors.ErrSlugNotFound
	}

	return slug, slugID, err
}

func (repo UrlDBRepository) GetFreeSlugID(slug string) (int, error) {
	slog.Debug("Getting last slug ID", slog.String("slug", slug))

	query := "SELECT COALESCE(MAX(slug_id), 0) as last_slug_id FROM urls where slug=$1"
	row := repo.db.QueryRow(query, slug)

	var slugID int

	err := row.Scan(&slugID)
	if err != nil {
		return 0, err
	}

	return slugID + 1, nil
}
