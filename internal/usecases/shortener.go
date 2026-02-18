package usecases

import (
	"urlshortener/internal/cache"
	"urlshortener/internal/database"
	"urlshortener/internal/errors"
)

type UrlUseCases struct {
	urlRepo  database.UrlRepository
	urlCache cache.UrlCache
}

func NewUrlUseCases(urlRepo database.UrlRepository, urlCache cache.UrlCache) UrlUseCases {
	return UrlUseCases{urlRepo, urlCache}
}

func (uc UrlUseCases) GetSlug(url string) (string, error) {
	slug, err := uc.urlRepo.GetSlugByUrl(url)
	if err != nil && err != errors.SlugNotFound {
		return "", err
	}

	if err == errors.SlugNotFound {
		slug = generateSlug(8)
		uc.urlRepo.Create(url, slug)
	}

	return slug, nil
}

func (uc UrlUseCases) GetUrl(slug string) (string, error) {
	return uc.urlRepo.GetUrlBySlug(slug)
}
