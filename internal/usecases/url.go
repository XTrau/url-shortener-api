package usecases

import (
	"errors"
	"math/rand"
	"strconv"
	"urlshortener/internal/apperrors"
	"urlshortener/internal/cache"
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
	urlRepository domain.UrlRepository
	urlCache      cache.UrlCacher
	g             *singleflight.Group
}

func NewUrlUseCases(urlRepo domain.UrlRepository, urlCache cache.UrlCacher) UrlUseCases {
	return UrlUseCases{
		urlRepository: urlRepo,
		urlCache:      urlCache,
		g:             &singleflight.Group{},
	}
}

// Сохраняет url в бд, создает для него новый slug
func (uc UrlUseCases) save(url domain.Url) (domain.Slug, error) {
	var slug domain.Slug

	// Проверка на случай если url в бд уже существует
	_, err := uc.urlRepository.GetSlugByUrl(url)

	if err == nil {
		return domain.Slug{}, apperrors.ErrAlreadyExists
	}

	if !errors.Is(err, apperrors.ErrSlugNotFound) {
		// Неизвестная ошибка
		return domain.Slug{}, err
	} else {
		// Ошибка Slug не найден, создаем новый и сохраняем в репозиторий

		slugText := generateSlug(8)
		slugID, err := uc.urlRepository.GetFreeSlugID(slugText)

		if err != nil {
			return domain.Slug{}, err
		}

		slug = domain.Slug{
			SlugID: slugID,
			Text:   slugText,
		}

		err = uc.urlRepository.Create(url, slug)
		if err != nil {
			return domain.Slug{}, err
		}

		return slug, nil
	}
}

// Возвращает slug для переданного url, если такого url не существует, то генерируется новый slug и сохраняется в бд
func (uc UrlUseCases) GetSlug(url domain.Url) (domain.Slug, error) {
	// Проверка в cache
	cachedSlug, err := uc.urlCache.GetSlug(url)
	if err == nil {
		return domain.Slug{
			Text:   cachedSlug.Text,
			SlugID: cachedSlug.SlugID,
		}, nil
	}

	// Проверка в repository, сохранение в cache, паттерн Single flight
	s, err, _ := uc.g.Do(url.Url, func() (interface{}, error) {
		slug, err := uc.urlRepository.GetSlugByUrl(url)

		if err != nil && !errors.Is(err, apperrors.ErrSlugNotFound) {
			return domain.Slug{}, err
		}

		if err != nil && errors.Is(err, apperrors.ErrSlugNotFound) {
			// Если slug для данного url не найден, то создаем новый в бд
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

// Возвращает url для переданного slug
func (uc UrlUseCases) GetUrl(slug domain.Slug) (domain.Url, error) {
	// Проверка в cache
	url, err := uc.urlCache.GetUrl(slug)
	if err == nil {
		return url, nil
	}

	// Group key для single flight
	groupKey := slug.Text + "_" + strconv.Itoa(slug.SlugID)

	// Получаем Slug с repository, сохранение в cache, паттерн Single flight
	u, err, _ := uc.g.Do(groupKey, func() (interface{}, error) {
		url, err = uc.urlRepository.GetUrlBySlug(slug)
		if err != nil {
			return domain.Url{}, err
		}
		err = uc.urlCache.Save(url, slug)
		return url, err
	})

	return u.(domain.Url), err
}
