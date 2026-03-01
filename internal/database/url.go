package database

import (
	"database/sql"
	"errors"
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
	_, err := repo.db.Exec("INSERT INTO urls (url, slug) VALUES ($1, $2)", url, slug)
	return err
}

func (repo UrlDBRepository) GetUrlBySlug(slug string) (string, error) {
	var url string
	row := repo.db.QueryRow("SELECT url FROM urls WHERE slug=$1", slug)
	err := row.Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return "", apperrors.ErrUrlNotFound
	}
	return url, err
}

func (repo UrlDBRepository) GetSlugByUrl(url string) (string, error) {
	var slug string
	row := repo.db.QueryRow("SELECT slug FROM urls WHERE url=$1", url)
	err := row.Scan(&slug)
	if errors.Is(err, sql.ErrNoRows) {
		return "", apperrors.ErrSlugNotFound
	}
	return slug, err
}
