package cache

import (
	"encoding/json"
	"os"
	"path"
	"time"
)

type fileCache struct {
	path string
}

// NewFileCache creates a new file cache
func NewFileCache(path string) Cache {
	return &fileCache{
		path: path,
	}
}

func (c *fileCache) Get(key string) (string, error) {
	f, err := os.Open(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	defer f.Close()

	var m map[string]cacheEntry

	err = json.NewDecoder(f).Decode(&m)
	if err != nil {
		return "", err
	}

	entry, ok := m[key]
	if !ok {
		return "", nil
	}

	if entry.Expires.Before(time.Now()) {
		return "", nil
	}

	return entry.Value, nil
}

func (c *fileCache) PutWithExpires(key string, value string, expires time.Time) error {
	err := os.MkdirAll(path.Dir(c.path), 0750)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(c.path, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		return err
	}

	defer f.Close()

	m := map[string]cacheEntry{}

	_ = json.NewDecoder(f).Decode(&m)

	m[key] = cacheEntry{
		Value:   value,
		Expires: expires,
	}

	f.Seek(0, 0)

	err = json.NewEncoder(f).Encode(m)
	if err != nil {
		return err
	}

	return nil
}
