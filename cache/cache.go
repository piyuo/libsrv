package file

import (
	"time"

	goCache "github.com/patrickmn/go-cache"
)

// Create a cache with a default expiration time of 5 minutes, and which
// purges expired items every 10 minutes
//
var cache = goCache.New(6*time.Minute, 1*time.Hour)

// Set Add an item to the cache, replacing any existing item. If the duration is 0
// (DefaultExpiration), the cache's default expiration time is used. If it is -1
// (NoExpiration), the item never expires.
//
func Set(key string, value interface{}, duration time.Duration) {
	cache.Set(key, value, duration)
}

// Get an item from the cache. Returns the item or nil, and a bool indicating whether the key was found.
//
func Get(key string) (interface{}, bool) {
	return cache.Get(key)
}
