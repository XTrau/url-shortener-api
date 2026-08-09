package database

import (
	"database/sql"
	"errors"
	"log/slog"
	"urlshortener/internal/apperrors"
	"urlshortener/internal/domain"
)

type UrlRepository interface {
	Create(url domain.Url, slug domain.Slug) error
	GetUrlBySlug(slug domain.Slug) (domain.Url, error)
	GetSlugByUrl(url domain.Url) (domain.Slug, error)
	GetFreeSlugID(slugText string) (int, error)
}

type UrlPostgresRepository struct {
	db *sql.DB
}

func NewUrlPostgresRepository(db *sql.DB) UrlPostgresRepository {
	return UrlPostgresRepository{db}
}

func (repo UrlPostgresRepository) Create(url domain.Url, slug domain.Slug) error {
	slog.Debug("Inserting url to database", slog.String("url", url.Url), slog.String("slug", slug.Text), slog.Int("slugID", slug.SlugID))

	query := "INSERT INTO urls (url, slug, slug_id) VALUES ($1, $2, $3)"
	_, err := repo.db.Exec(query, url, slug.Text, slug.SlugID)

	return err
}

func (repo UrlPostgresRepository) GetUrlBySlug(slug domain.Slug) (domain.Url, error) {
	slog.Debug("Getting url from database", slog.String("slug", slug.Text), slog.Int("slugID", slug.SlugID))

	query := "SELECT url FROM urls WHERE slug=$1 AND slug_id=$2"
	row := repo.db.QueryRow(query, slug.Text, slug.SlugID)

	var url string
	err := row.Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Url{}, apperrors.ErrUrlNotFound
	}
	return domain.Url{Url: url}, err
}

func (repo UrlPostgresRepository) GetSlugByUrl(url domain.Url) (domain.Slug, error) {
	slog.Debug("Getting url from database", slog.String("url", url.Url))

	query := "SELECT slug, slug_id FROM urls WHERE url=$1"
	row := repo.db.QueryRow(query, url)

	var slug string
	var slugID int

	err := row.Scan(&slug, &slugID)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Slug{}, apperrors.ErrSlugNotFound
	}

	return domain.Slug{
		Text:   slug,
		SlugID: slugID,
	}, err
}

func (repo UrlPostgresRepository) GetFreeSlugID(slugText string) (int, error) {
	slog.Debug("Getting last slug ID", slog.String("slug", slugText))

	query := "SELECT COALESCE(MAX(slug_id), 0) as last_slug_id FROM urls where slug=$1"
	row := repo.db.QueryRow(query, slugText)

	var slugID int

	err := row.Scan(&slugID)
	if err != nil {
		return 0, err
	}

	return slugID + 1, nil
}
