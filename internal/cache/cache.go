package cache

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
	"urlshortener/internal/domain"
)

const UrlTTL = time.Minute * 5

type UrlCacher interface {
	Save(url domain.Url, slug domain.Slug) error
	GetUrl(slug domain.Slug) (domain.Url, error)
	GetSlug(url domain.Url) (domain.Slug, error)
}

type UrlCache struct {
	cache domain.Cache
}

func NewUrlCache(storage domain.Cache) UrlCache {
	return UrlCache{storage}
}

func (uc UrlCache) GetUrlKey(slug string, slugID int) string {
	return fmt.Sprintf("url:%s#%d", slug, slugID)
}

func (uc UrlCache) GetSlugKey(url string) string {
	return fmt.Sprintf("slug:%s", url)
}

func (uc UrlCache) Save(url domain.Url, slug domain.Slug) error {
	urlKey := uc.GetUrlKey(slug.Text, slug.SlugID)
	slugKey := uc.GetSlugKey(url.Url)

	slog.Debug(
		"Saving url to Cache",
		slog.String("urlKey", urlKey),
		slog.String("slugKey", slugKey),
	)

	slugData, err := json.Marshal(slug)

	if err != nil {
		return err
	}

	err = uc.cache.Set(slugKey, string(slugData), UrlTTL)
	if err != nil {
		return err
	}

	err = uc.cache.Set(urlKey, url.Url, UrlTTL)
	if err != nil {
		return err
	}

	return nil
}

func (uc UrlCache) GetUrl(slug domain.Slug) (domain.Url, error) {
	key := uc.GetUrlKey(slug.Text, slug.SlugID)
	slog.Debug("Getting url from Cache", slog.String("key", key))

	url, err := uc.cache.Get(key)

	if err != nil {
		return domain.Url{}, err
	}

	return domain.Url{Url: url}, err
}

func (uc UrlCache) GetSlug(url domain.Url) (domain.Slug, error) {
	key := uc.GetSlugKey(url.Url)
	slog.Debug("Getting slug from Redis", slog.String("key", key))

	slugData, err := uc.cache.Get(key)

	if err != nil {
		return domain.Slug{}, err
	}

	var cachedSlug domain.Slug
	if err := json.Unmarshal([]byte(slugData), &cachedSlug); err != nil {
		return domain.Slug{}, err
	}

	return cachedSlug, nil
}
