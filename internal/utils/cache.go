package utils

import (
	"os"
	"sync"
)

type DownloadCache struct {
	data map[string]string
	mu   sync.RWMutex
}

var FileCache = &DownloadCache{
	data: make(map[string]string),
}

func (c *DownloadCache) Get(key string) (string, bool) {
	c.mu.RLock()
	filePath, exists := c.data[key]
	c.mu.RUnlock()

	if !exists {
		return "", false
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.Delete(key)
		return "", false
	}
	return filePath, true
}

func (c *DownloadCache) Set(key, filePath string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = filePath
}

func (c *DownloadCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}
