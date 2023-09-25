package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	defaultExpiration      = time.Minute * 5
	defaultCleanUpDuration = time.Minute * 8
)

type Cache struct {
	cache *cache.Cache
}

func NewCache() *Cache {
	return &Cache{
		cache: cache.New(defaultExpiration, defaultCleanUpDuration),
	}
}
