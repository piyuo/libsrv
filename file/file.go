package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/piyuo/libsrv/cache"
	"github.com/pkg/errors"
)

const (
	CacheKey = "f-"
)

var baseDir string = ""

type CacheType int8

const (
	NoCache CacheType = iota
	Cache
	GzipCache
)

// Read binary data from file, filename can be relative path
//
//	bytes, err := file.Read("mock/mock.json", NoCache)
//
func Read(filename string, cacheType CacheType, d time.Duration) ([]byte, error) {
	// check cache
	cacheKey := CacheKey + filename

	switch cacheType {
	case Cache:
		found, bytes, err := cache.Get(cacheKey)
		if err != nil {
			return nil, errors.Wrap(err, "get cache "+cacheKey)
		}
		if found {
			return bytes, nil
		}
	case GzipCache:
		found, bytes, err := cache.GzipGet(cacheKey)
		if err != nil {
			return nil, errors.Wrap(err, "get gzip cache "+cacheKey)
		}
		if found {
			return bytes, nil
		}
	}

	// read file
	fullPath := Lookup(filename)
	if fullPath == "" {
		return nil, errors.New(filename + " not found")
	}

	osFile, err := os.Open(fullPath)
	if err != nil {
		return nil, errors.Wrapf(err, "open %v", fullPath)
	}

	content, err := ioutil.ReadAll(osFile)
	if err != nil {
		return nil, errors.Wrapf(err, "read %v", fullPath)
	}

	if err := osFile.Close(); err != nil {
		return nil, errors.Wrapf(err, "close %v", fullPath)
	}

	// write cache
	switch cacheType {
	case Cache:
		if err := cache.Set(cacheKey, content, d); err != nil {
			return nil, errors.Wrap(err, "set cache "+cacheKey)
		}
	case GzipCache:
		if err := cache.GzipSet(cacheKey, content, d); err != nil {
			return nil, errors.Wrap(err, "set cache "+cacheKey)
		}
	}
	return content, nil
}

// ReadText read text from file, filename can be relative path
//
//	txt, err := file.ReadText("mock/mock.json")
//
func ReadText(filename string, cacheType CacheType, d time.Duration) (string, error) {
	bytes, err := Read(filename, cacheType, d)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ReadJSON read json object from file, filename can be relative path
//
//	f, err := file.ReadJSON("mock/mock.json")
//
func ReadJSON(filename string, cacheType CacheType, d time.Duration) (map[string]interface{}, error) {
	bytes, err := Read(filename, cacheType, d)
	if err != nil {
		return nil, err
	}
	content := make(map[string]interface{})
	if err := json.Unmarshal([]byte(bytes), &content); err != nil {
		return nil, errors.Wrapf(err, "decode json %v", filename)
	}
	return content, nil
}

// Exists return true if file found
//
//	f, err := file.ReadJSON("mock/mock.json")
//
func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		return false
	}
	return true
}

// Lookup find dir or file from current path all the way to the top, return actual path where dir or file locate, return empty if can not determine
//
//	fullPath := Lookup("keys")
//
func Lookup(name string) string {
	if baseDir != "" {
		return path.Join(baseDir, name)
	}

	curdir, err := os.Getwd()
	if err != nil {
		return ""
	}

	var filepath string
	dir := curdir
	for i := 0; i <= 5; i++ {
		filepath = path.Join(dir, name)
		if Exists(filepath) {
			// target found. we know the base now
			fmt.Printf("base dir is %v \n", dir)
			baseDir = dir
			return filepath
		}

		// already root dir, just give up
		if dir == "/" {
			break
		}
		//dir not exist, go up
		dir = path.Join(dir, "../")
	}
	//	fmt.Printf("%v not found in %v or parent dir\n", name, curdir)
	return ""
}
