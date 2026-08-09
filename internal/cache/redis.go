package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

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
