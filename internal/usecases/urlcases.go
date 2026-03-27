package usecases

import (
	"errors"
	"urlshortener/internal/apperrors"
	"urlshortener/internal/cache"
	"urlshortener/internal/database"
)

type UrlUseCases struct {
	urlRepo  database.UrlRepository
	urlCache cache.UrlCacher
}

func NewUrlUseCases(urlRepo database.UrlRepository, urlCache cache.UrlCacher) UrlUseCases {
	return UrlUseCases{urlRepo, urlCache}
}

func (uc UrlUseCases) GetSlug(url string) (string, int, error) {
	cachedSlug, err := uc.urlCache.GetSlug(url)
	if err == nil {
		return cachedSlug.Slug, cachedSlug.SlugID, nil
	}

	slug, slugID, err := uc.urlRepo.GetSlugByUrl(url)
	if err != nil && !errors.Is(err, apperrors.ErrSlugNotFound) {
		return "", 0, err
	}

	if errors.Is(err, apperrors.ErrSlugNotFound) {
		slug = generateSlug(8)
		slugID, err = uc.urlRepo.GetFreeSlugID(slug)

		if err != nil {
			return "", 0, err
		}

		err := uc.urlRepo.Create(url, slug, slugID)

		if err != nil {
			return "", 0, err
		}
	}

	err = uc.urlCache.Save(url, slug, slugID)

	return slug, slugID, nil
}

func (uc UrlUseCases) GetUrl(slug string, slugID int) (string, error) {
	url, err := uc.urlCache.GetUrl(slug, slugID)
	if err == nil {
		return url, nil
	}

	url, err = uc.urlRepo.GetUrlBySlug(slug, slugID)

	if err != nil {
		return "", err
	}

	err = uc.urlCache.Save(url, slug, slugID)

	return url, nil
}
