package usecases

import (
	"errors"
	"math/rand"
	"strconv"
	"urlshortener/internal/apperrors"
	"urlshortener/internal/cache"
	"urlshortener/internal/database"
	"urlshortener/internal/domain"

	"golang.org/x/sync/singleflight"
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
	g        *singleflight.Group
}

func NewUrlUseCases(urlRepo database.UrlRepository, urlCache cache.UrlCacher) UrlUseCases {
	return UrlUseCases{
		urlRepo:  urlRepo,
		urlCache: urlCache,
		g:        &singleflight.Group{},
	}
}

func (uc UrlUseCases) save(url domain.Url) (domain.Slug, error) {
	var slug domain.Slug
	_, err := uc.urlRepo.GetSlugByUrl(url)

	if err == nil {
		return domain.Slug{}, apperrors.ErrAlreadyExists
	}

	if !errors.Is(err, apperrors.ErrSlugNotFound) {
		return domain.Slug{}, err
	} else {
		slugText := generateSlug(8)
		slugID, err := uc.urlRepo.GetFreeSlugID(slugText)

		if err != nil {
			return domain.Slug{}, err
		}

		slug = domain.Slug{
			SlugID: slugID,
			Text:   slugText,
		}

		err = uc.urlRepo.Create(url, slug)
		if err != nil {
			return domain.Slug{}, err
		}

		return slug, nil
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

	groupKey := url.Url

	s, err, _ := uc.g.Do(groupKey, func() (interface{}, error) {
		slug, err := uc.urlRepo.GetSlugByUrl(url)

		if err != nil && !errors.Is(err, apperrors.ErrSlugNotFound) {
			return domain.Slug{}, err
		}

		if err != nil && errors.Is(err, apperrors.ErrSlugNotFound) {
			slug, err = uc.save(url)

			if err != nil {
				return domain.Slug{}, err
			}
		}

		uc.urlCache.Save(url, slug)
		return slug, err
	})

	return s.(domain.Slug), err
}

func (uc UrlUseCases) GetUrl(slug domain.Slug) (domain.Url, error) {
	url, err := uc.urlCache.GetUrl(slug)
	if err == nil {
		return url, nil
	}

	groupKey := slug.Text + "_" + strconv.Itoa(slug.SlugID)

	u, err, _ := uc.g.Do(groupKey, func() (interface{}, error) {
		url, err = uc.urlRepo.GetUrlBySlug(slug)
		if err != nil {
			return domain.Url{}, err
		}
		err = uc.urlCache.Save(url, slug)
		return url, err
	})

	return u.(domain.Url), err
}
