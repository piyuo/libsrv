package cache

import (
	"time"

	goCache "github.com/patrickmn/go-cache"
)

// Frequency define cache item use frequency
//
type Frequency int

const (
	// LOW frequency only keep item in 3 min and cached item total count must less than 1000
	//
	LOW Frequency = iota

	// MEDIUM frequency Frequency keep item in 10 min and cached item total count must less than 2,000
	//
	MEDIUM

	// HIGH frequency Frequency will keep in cache 24 hour and total count must less than 2500
	//
	HIGH
)

// cache a cache with a default expiration time of 5 minutes, and which purges expired items every 4 minutes
//
var cache = goCache.New(10*time.Minute, 3*time.Minute)

// Set Add an item to the cache, replacing any existing item. you need specify frequency this item will be use,
//
// Low Frequency = cache 3 min, total limit 1,000
//
// Medium Frequency = cache 10 min, total limit 2,000
//
// High Frequency = cache 24 hour , total limit 2,500
//
func Set(freq Frequency, key string, value interface{}) {
	var d time.Duration
	switch freq {
	case LOW:
		d = 180 * time.Second
		if Count() >= 1000 {
			return
		}
	case MEDIUM:
		d = 600 * time.Second
		if Count() >= 2000 {
			return
		}
	default:
		d = 24 * time.Hour
		if Count() >= 2500 {
			return
		}
	}
	cache.Set(key, value, d)
}

// Increment a value in cache, set value if cache not exist
//
//	Increment(LOW, key, 1)
//
func Increment(freq Frequency, key string, value int) {
	record, found := cache.Get(key)
	if !found {
		Set(freq, key, value)
		return
	}
	count := record.(int)
	count = count + value
	Set(freq, key, count)
}

func set(key string, value interface{}, duration time.Duration) {
	cache.Set(key, value, duration)
}

// Get an item from the cache. Returns the item or nil, and a bool indicating whether the key was found.
//
func Get(key string) (interface{}, bool) {
	return cache.Get(key)
}

// GetBool returns the bool and a bool indicating whether the key was found.
//
func GetBool(key string) (bool, bool) {
	item, found := cache.Get(key)
	if !found {
		return false, false
	}
	return item.(bool), found
}

// GetInt returns the int and a bool indicating whether the key was found.
//
func GetInt(key string) (int, bool) {
	item, found := cache.Get(key)
	if !found {
		return 0, false
	}
	return item.(int), found
}

// GetInt64 returns the int and a bool indicating whether the key was found.
//
func GetInt64(key string) (int64, bool) {
	item, found := cache.Get(key)
	if !found {
		return 0, false
	}
	return item.(int64), found
}

// GetUInt32 returns the uint32 and a bool indicating whether the key was found.
//
func GetUInt32(key string) (uint32, bool) {
	item, found := cache.Get(key)
	if !found {
		return 0, false
	}
	return item.(uint32), found
}

// GetUInt64 returns the uint64 and a bool indicating whether the key was found.
//
func GetUInt64(key string) (uint64, bool) {
	item, found := cache.Get(key)
	if !found {
		return 0, false
	}
	return item.(uint64), found
}

// GetString returns the string or nil, and a bool indicating whether the key was found.
//
func GetString(key string) (string, bool) {
	item, found := cache.Get(key)
	if !found {
		return "", false
	}
	return item.(string), found
}

// GetBytes returns the bytes or nil, and a bool indicating whether the key was found.
//
func GetBytes(key string) ([]byte, bool) {
	item, found := cache.Get(key)
	if !found {
		return nil, false
	}
	return item.([]byte), found
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

// Count return cache item total count
//
func Count() int {
	return cache.ItemCount()
}
