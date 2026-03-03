package cache

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"urlshortener/internal/apperrors"

	"github.com/redis/go-redis/v9"
)

type UrlCache interface {
	Save(url string, slug string) error
	GetUrl(slug string) (string, error)
	GetSlug(url string) (string, error)
}

type UrlRedisCache struct {
	rdb *redis.Client
}

func NewUrlRedisCache(rdb *redis.Client) UrlRedisCache {
	return UrlRedisCache{rdb}
}

func (uc UrlRedisCache) GetUrlKey(slug string) string {
	return fmt.Sprintf("url:%s", slug)
}

func (uc UrlRedisCache) GetSlugKey(url string) string {
	return fmt.Sprintf("slug:%s", url)
}

func (uc UrlRedisCache) Save(url string, slug string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	urlKey := uc.GetUrlKey(slug)
	slugKey := uc.GetSlugKey(url)

	slog.Debug(
		"Saving url to Redis",
		slog.String("urlKey", urlKey),
		slog.String("slugKey", slugKey),
	)

	err := uc.rdb.Set(ctx, slugKey, slug, time.Minute*5).Err()

	if err != nil {
		return err
	}

	err = uc.rdb.Set(ctx, urlKey, url, time.Minute*5).Err()

	return err
}

func (uc UrlRedisCache) GetUrl(slug string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	key := uc.GetUrlKey(slug)

	slog.Debug("Getting url from Redis", slog.String("key", key))

	url, err := uc.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return "", apperrors.ErrCacheKeyNotFound
	}

	return url, err
}

func (uc UrlRedisCache) GetSlug(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	key := uc.GetSlugKey(url)

	slog.Debug("Getting slug from Redis", slog.String("key", key))

	url, err := uc.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return "", apperrors.ErrCacheKeyNotFound
	}

	return url, err
}
