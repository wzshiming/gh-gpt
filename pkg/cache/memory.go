package cache

import (
	"time"
)

type memoryCache struct {
	m map[string]cacheEntry
}

// NewMemoryCache creates a new memory cache
func NewMemoryCache() Cache {
	return &memoryCache{
		m: map[string]cacheEntry{},
	}
}

func (c *memoryCache) Get(key string) (string, error) {
	entry, ok := c.m[key]
	if !ok {
		return "", nil
	}

	if entry.Expires.Before(time.Now()) {
		return "", nil
	}

	return entry.Value, nil
}

func (c *memoryCache) PutWithExpires(key string, value string, expires time.Time) error {
	c.m[key] = cacheEntry{
		Value:   value,
		Expires: expires,
	}

	return nil
}
