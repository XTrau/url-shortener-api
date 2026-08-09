package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig interface {
	RedisHost() string
	RedisPort() string
	RedisUser() string
	RedisPassword() string
	RedisDatabase() int
}

func NewRedisClient(cfg RedisConfig) (*redis.Client, error) {
	addr := fmt.Sprintf("%v:%v", cfg.RedisHost(), cfg.RedisPort())

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: cfg.RedisUser(),
		Password: cfg.RedisPassword(),
		DB:       cfg.RedisDatabase(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
