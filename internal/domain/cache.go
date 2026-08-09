package domain

import "time"

type KeyValueStorage interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration) error
}
