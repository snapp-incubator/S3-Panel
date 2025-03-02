package cache

import (
	"fmt"
	"sync"
)

type InMemoryCache struct {
	DB *sync.Map
}

func NewInMemoryCache(db *sync.Map) *InMemoryCache {
	return &InMemoryCache{DB: db}
}

func (i InMemoryCache) Set(targetKey, targetValue string) error {
	i.DB.Store(targetKey, targetValue)
	return nil
}

func (i InMemoryCache) Get(targetKey string) (string, error) {
	if value, ok := i.DB.Load(targetKey); ok {
		return value.(string), nil
	}
	return "", fmt.Errorf("key not found")
}
