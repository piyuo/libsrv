package file

import (
	"time"

	goCache "github.com/patrickmn/go-cache"
)

// Create a cache with a default expiration time of 5 minutes, and which
// purges expired items every 10 minutes
//
var cache = goCache.New(5*time.Minute, 10*time.Minute)

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

// Reset clear entire cache
//
func Reset() {
	cache.Flush()
}

// Delete delete item from cache
//
func Delete(key string) {
	cache.Delete(key)
}

// configCache set new cache with a given default expiration duration and cleanup interval. If the expiration duration is less than one (or NoExpiration), the items in the cache never expire (by default), and must be deleted manually. If the cleanup interval is less than one, expired items are not deleted from the cache before calling c.DeleteExpired().
//
func configCache(defaultExpiration, cleanupInterval time.Duration) {
	cache = goCache.New(defaultExpiration, cleanupInterval)
}
