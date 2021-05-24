package cache

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/coocood/freecache"
	"github.com/piyuo/libsrv/digit"
	"github.com/pkg/errors"
)

// In bytes, where 1024 * 1024 represents a single Megabyte, and 20 * 1024*1024 represents 20 Megabytes.
var cache = freecache.NewCache(20 * 1024 * 1024)

// defaultDuration is 20 minutes
const defaultDuration = 20 * time.Minute

// Set an entry to the cache, replacing any existing entry. If the duration is 0 (DefaultExpiration), the cache's default 20 minutes, if duration is -1 there is no cache
//
//	err = Set("key1", []byte("hi"), 0)
//
func Set(key string, value []byte, d time.Duration) error {
	if d == 0 {
		d = defaultDuration
	}
	if d == -1 {
		return nil
	}

	fmt.Printf("%v cached\n", key)
	return cache.Set([]byte(key), value, int(d.Seconds()))
}

// Set an entry to the cache, replacing any existing entry. If the duration is 0 (DefaultExpiration), the cache's default 20 minutes is used
//
//	err = SetString(123, "hi", 0)
//
func ISet(key int64, value []byte, d time.Duration) error {
	if d == 0 {
		d = defaultDuration
	}
	if d == -1 {
		return nil
	}
	return cache.SetInt(key, value, int(d.Seconds()))
}

// Get an entry from the cache. Returns the item or nil, and a bool indicating whether the key was found.
//
//	found, bytes, err := Get("key1")
//
func Get(key string) (bool, []byte, error) {
	value, err := cache.Get([]byte(key))
	if err != nil {
		if err == freecache.ErrNotFound {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, value, nil
}

// Get an entry from the cache. Returns the item or nil, and a bool indicating whether the key was found.
//
//	found, bytes, err := IGet(123)
//
func IGet(key int64) (bool, []byte, error) {
	value, err := cache.GetInt(key)
	if err != nil {
		if err == freecache.ErrNotFound {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, value, nil
}

// SetString set string entry to the cache, replacing any existing entry. If the duration is 0 (DefaultExpiration), the cache's default 20 minutes is used
//
//	err := SetString(key, "hi", 0)
//
func SetString(key, value string, d time.Duration) error {
	return Set(key, []byte(value), d)
}

// GetString return an string from the cache. Returns the item or nil, and a bool indicating whether the key was found.
//
//	found, str, err := GetString(key)
//
func GetString(key string) (bool, string, error) {
	found, value, err := Get(key)
	if err != nil {
		return false, "", err
	}
	if !found {
		return false, "", nil
	}
	return true, string(value), nil
}

// SetInt64 set int64 entry to the cache, replacing any existing entry. If the duration is 0 (DefaultExpiration), the cache's default 20 minutes is used
//
//	err := SetInt64(key, 65536, 0)
//
func SetInt64(key string, value int64, d time.Duration) error {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(value))
	return Set(key, bytes, d)
}

// GetInt64 return an int64 from the cache. Returns the item or nil, and a bool indicating whether the key was found.
//
//	found, valueInt64, err := GetInt64(key)
//
func GetInt64(key string) (bool, int64, error) {
	found, bytes, err := Get(key)
	if err != nil {
		return false, -1, err
	}
	if !found {
		return false, -1, nil
	}
	i := int64(binary.LittleEndian.Uint64(bytes))
	return true, i, nil
}

// Delete delete entry from cache
//
//	Delete("key1")
//
func Delete(key string) {
	cache.Del([]byte(key))
}

// IDelete delete entry from cache
//
//	IDelete(123)
//
func IDelete(key int64) {
	cache.DelInt(key)
}

// Count return cache entry total count
//
//	count := Count()
//
func Count() int64 {
	return cache.EntryCount()
}

// GzipSet compress byte array before set an entry to the cache
//
//	err = GzipSet("key1", []byte("hi"), 0)
//
func GzipSet(key string, value []byte, d time.Duration) error {
	zipped, err := digit.Compress(value)
	if err != nil {
		return errors.Wrap(err, "compress")
	}
	return Set(key, zipped, d)
}

// GzipGet get an entry from the cache and decompress zipped content
//
//	found, bytes, err := GzipGet("key1")
//
func GzipGet(key string) (bool, []byte, error) {
	found, value, err := Get(key)
	if err != nil {
		return false, nil, errors.Wrap(err, "get")
	}
	unzipped, err := digit.Decompress(value)
	if err != nil {
		return false, nil, errors.Wrap(err, "decompress")
	}
	return found, unzipped, nil
}
