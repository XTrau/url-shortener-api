package cache_test

import (
	"sync"
	"testing"
	"time"
	"urlshortener/internal/cache"
	"urlshortener/internal/domain"
)

// Struct without ttl for tests
type MapKeyValueStorage struct {
	cache map[string]string
	mu    sync.RWMutex
}

func (kv *MapKeyValueStorage) Get(key string) (string, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	value, ok := kv.cache[key]
	if !ok {
		return "", cache.ErrCacheKeyNotFound
	}
	return value, nil
}

func (kv *MapKeyValueStorage) Set(key string, value string, ttl time.Duration) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.cache[key] = value
	return nil
}

func TestUrlCache(t *testing.T) {
	cacheStore := &MapKeyValueStorage{
		cache: make(map[string]string),
		mu:    sync.RWMutex{},
	}
	urlCache := cache.NewUrlCache(cacheStore)

	url := domain.Url{
		Url: "http://example.com",
	}

	slug := domain.Slug{
		Text:   "qweqwe",
		SlugID: 1,
	}

	err := urlCache.Save(url, slug)
	if err != nil {
		t.Errorf("urlCache.Save must return nil error")
	}

	cachedSlug, err := urlCache.GetSlug(url)

	if err != nil {
		t.Errorf("urlCache.GetSlug must return nil error")
	}

	if cachedSlug.Text != slug.Text || cachedSlug.SlugID != slug.SlugID {
		t.Errorf(
			"cachedSlug.Slug=%q, cachedSlug.SlugID=%d expected: cachedSlug.Slug=%q, cachedSlug.SlugID=%d",
			cachedSlug.Text,
			cachedSlug.SlugID,
			slug.Text,
			slug.SlugID,
		)
	}

	cachedURL, err := urlCache.GetUrl(slug)

	if err != nil {
		t.Errorf("urlCache.GetUrl must return nil error")
	}

	if cachedURL != url {
		t.Errorf("cachedURL=%q expected: %q", cachedURL, url)
	}
}
