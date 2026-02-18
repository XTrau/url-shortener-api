package usecases

import (
	"urlshortener/internal/database"
	"urlshortener/internal/errors"
)

type UrlUseCases struct {
	repo database.UrlRepository
}

func NewUrlUseCases(repo database.UrlRepository) UrlUseCases {
	return UrlUseCases{repo}
}

func (uc UrlUseCases) GetSlug(url string) (string, error) {
	slug, err := uc.repo.GetSlugByUrl(url)
	if err != nil && err != errors.SlugNotFound {
		return "", err
	}

	if err == errors.SlugNotFound {
		slug = generateSlug(8)
		uc.repo.Create(url, slug)
	}

	return slug, nil
}

func (uc UrlUseCases) GetUrl(slug string) (string, error) {
	return uc.repo.GetUrlBySlug(slug)
}
