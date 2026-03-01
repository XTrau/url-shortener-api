package usecases

import (
	"errors"
	"log"
	"urlshortener/internal/apperrors"
	"urlshortener/internal/cache"
	"urlshortener/internal/database"
)

type UrlUseCases struct {
	urlRepo  database.UrlRepository
	urlCache cache.UrlCache
}

func NewUrlUseCases(urlRepo database.UrlRepository, urlCache cache.UrlCache) UrlUseCases {
	return UrlUseCases{urlRepo, urlCache}
}

func (uc UrlUseCases) GetSlug(url string) (string, error) {
	slug, err := uc.urlCache.GetSlug(url)
	if err == nil {
		return slug, nil
	}

	if !errors.Is(err, apperrors.ErrCacheKeyNotFound) {
		return "", err
	}

	slug, err = uc.urlRepo.GetSlugByUrl(url)
	if err != nil && !errors.Is(err, apperrors.ErrSlugNotFound) {
		return "", err
	}

	if errors.Is(err, apperrors.ErrSlugNotFound) {
		slug = generateSlug(8)
		uc.urlRepo.Create(url, slug)
	}

	err = uc.urlCache.Save(url, slug)
	if err != nil {
		log.Println("Error on save to Redis.", err)
	}

	return slug, nil
}

func (uc UrlUseCases) GetUrl(slug string) (string, error) {
	url, err := uc.urlCache.GetUrl(slug)
	if err == nil {
		return url, nil
	}

	if !errors.Is(err, apperrors.ErrCacheKeyNotFound) {
		return "", err
	}

	url, err = uc.urlRepo.GetUrlBySlug(slug)

	err = uc.urlCache.Save(url, slug)
	if err != nil {
		log.Println("Error on save to Redis.", err)
	}

	return url, nil
}
