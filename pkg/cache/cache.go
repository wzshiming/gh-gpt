package cache

import (
	"time"
)

// Cache is a cache interface
type Cache interface {
	PutWithExpires(key string, value string, expires time.Time) error
	Get(key string) (string, error)
}

type cacheEntry struct {
	Value   string    `json:"value"`
	Expires time.Time `json:"expires"`
}
