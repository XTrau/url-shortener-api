package cache

import (
	"github.com/redis/go-redis/v9"
)

type UrlCache interface {
	GetUrl(slug string) (string, error)
	GetSlug(url string) (string, error)
}

type UrlRedisCache struct {
	rdb *redis.Client
}

func NewUrlRedisCache(rdb *redis.Client) UrlRedisCache {
	return UrlRedisCache{rdb}
}

func (uc UrlRedisCache) GetUrl(slug string) (string, error) {
	var result string
	return result, nil
}

func (uc UrlRedisCache) GetSlug(url string) (string, error) {
	var result string
	return result, nil
}
