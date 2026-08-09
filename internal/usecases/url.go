package usecases

import (
	"errors"
	"math/rand"
	"urlshortener/internal/apperrors"
	"urlshortener/internal/cache"
	"urlshortener/internal/database"
	"urlshortener/internal/domain"
)

const charSet = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM134567890"

func generateSlug(size int) string {
	var result string
	for range size {
		ind := rand.Intn(len(charSet))
		result += charSet[ind : ind+1]
	}
	return result
}

type UrlUseCases struct {
	urlRepo  database.UrlRepository
	urlCache cache.UrlCacher
}

func NewUrlUseCases(urlRepo database.UrlRepository, urlCache cache.UrlCacher) UrlUseCases {
	return UrlUseCases{
		urlRepo:  urlRepo,
		urlCache: urlCache,
	}
}

func (uc UrlUseCases) GetSlug(url domain.Url) (domain.Slug, error) {
	cachedSlug, err := uc.urlCache.GetSlug(url)
	if err == nil {
		return domain.Slug{
			Text:   cachedSlug.Text,
			SlugID: cachedSlug.SlugID,
		}, nil
	}

	slug, err := uc.urlRepo.GetSlugByUrl(url)
	if err != nil && !errors.Is(err, apperrors.ErrSlugNotFound) {
		return domain.Slug{}, err
	}

	if errors.Is(err, apperrors.ErrSlugNotFound) {
		slugText := generateSlug(8)
		slugID, err := uc.urlRepo.GetFreeSlugID(slugText)

		if err != nil {
			return domain.Slug{}, err
		}

		slug := domain.Slug{
			SlugID: slugID,
			Text:   slugText,
		}

		err = uc.urlRepo.Create(url, slug)

		if err != nil {
			return domain.Slug{}, err
		}
	}

	err = uc.urlCache.Save(url, slug)

	return slug, nil
}

func (uc UrlUseCases) GetUrl(slug domain.Slug) (domain.Url, error) {
	url, err := uc.urlCache.GetUrl(slug)
	if err == nil {
		return url, nil
	}

	url, err = uc.urlRepo.GetUrlBySlug(slug)
	if err != nil {
		return domain.Url{}, err
	}

	err = uc.urlCache.Save(url, slug)
	return url, nil
}
