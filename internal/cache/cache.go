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

func (uc UrlRedisCache) GetSlugKey(slug string) string {
	return fmt.Sprintf("slug:%s", slug)
}

func (uc UrlRedisCache) GetUrlKey(url string) string {
	return fmt.Sprintf("url:%s", url)
}

func (uc UrlRedisCache) Save(url string, slug string) error {
	slog.Debug("Saving url to Redis", slog.String("url", url), slog.String("slug", slug))

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	urlKey := uc.GetUrlKey(url)
	slugKey := uc.GetSlugKey(slug)

	err := uc.rdb.Set(ctx, urlKey, slug, time.Minute*5).Err()

	if err != nil {
		return err
	}

	err = uc.rdb.Set(ctx, slugKey, url, time.Minute*5).Err()

	return err
}

func (uc UrlRedisCache) GetUrl(slug string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	key := uc.GetSlugKey(slug)

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

	key := uc.GetUrlKey(url)

	slog.Debug("Getting url from Redis", slog.String("key", key))

	url, err := uc.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return "", apperrors.ErrCacheKeyNotFound
	}

	return url, err
}
