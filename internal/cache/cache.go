package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type CachedSlug struct {
	Slug   string `json:"slug"`
	SlugID int    `json:"slug_id"`
}

type KeyValueStorage interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration) error
}

type RedisKeyValueStorage struct {
	rdb     *redis.Client
	timeout time.Duration
}

func NewRedisKeyValueStorage(rdb *redis.Client, timeout time.Duration) RedisKeyValueStorage {
	return RedisKeyValueStorage{
		rdb:     rdb,
		timeout: timeout,
	}
}

func (kvs RedisKeyValueStorage) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kvs.timeout)
	defer cancel()

	res, err := kvs.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return "", ErrCacheKeyNotFound
	}

	return res, err
}

func (kvs RedisKeyValueStorage) Set(key string, value string, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), kvs.timeout)
	defer cancel()
	return kvs.rdb.Set(ctx, key, value, ttl).Err()
}

type UrlCacher interface {
	Save(url string, slug string, slugID int) error
	GetUrl(slug string, slugID int) (string, error)
	GetSlug(url string) (CachedSlug, error)
}

type UrlCache struct {
	cache KeyValueStorage
}

func NewUrlCache(storage KeyValueStorage) UrlCache {
	return UrlCache{storage}
}

func (uc UrlCache) GetUrlKey(slug string, slugID int) string {
	return fmt.Sprintf("url:%s#%d", slug, slugID)
}

func (uc UrlCache) GetSlugKey(url string) string {
	return fmt.Sprintf("slug:%s", url)
}

func (uc UrlCache) Save(url string, slug string, slugID int) error {
	const TTL = time.Minute * 5

	urlKey := uc.GetUrlKey(slug, slugID)
	slugKey := uc.GetSlugKey(url)

	slog.Debug(
		"Saving url to Cache",
		slog.String("urlKey", urlKey),
		slog.String("slugKey", slugKey),
	)

	cacheSlug := CachedSlug{slug, slugID}
	cachedSlugData, err := json.Marshal(cacheSlug)

	if err != nil {
		return err
	}

	err = uc.cache.Set(slugKey, string(cachedSlugData), TTL)

	if err != nil {
		return err
	}

	err = uc.cache.Set(urlKey, url, TTL)

	return err
}

func (uc UrlCache) GetUrl(slug string, slugID int) (string, error) {
	key := uc.GetUrlKey(slug, slugID)

	slog.Debug("Getting url from Cache", slog.String("key", key))

	url, err := uc.cache.Get(key)

	if err != nil {
		return "", err
	}

	return url, err
}

func (uc UrlCache) GetSlug(url string) (CachedSlug, error) {
	key := uc.GetSlugKey(url)

	slog.Debug("Getting slug from Redis", slog.String("key", key))

	cachedSlugStr, err := uc.cache.Get(key)

	if err != nil {
		return CachedSlug{}, err
	}

	var cachedSlug CachedSlug
	err = json.Unmarshal([]byte(cachedSlugStr), &cachedSlug)

	if err != nil {
		return CachedSlug{}, err
	}

	return cachedSlug, nil
}
