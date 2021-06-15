package file

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/piyuo/libsrv/cache"
	"github.com/pkg/errors"
)

const (
	CacheKey  = "f-"
	KeysDir   = "keys"
	AssetsDir = "assets"
	I18nDir   = "i18n"
)

// assetsBase hold assets base dir for lookup
var assetsBase string = ""

// keyBase hold key base dir for lookup
var keysBase string = ""

type CacheType int8

const (
	NoCache CacheType = iota
	Cache
	GzipCache
)

// Read binary data from file, filename can be relative path, It will not have error if file not found or can't be read, return nil instead
//
//	bytes, err := Read(AssetsDir,"i18n/mock_en_US.json", NoCache, 0)
//
func Read(baseDir, filename string, cacheType CacheType, d time.Duration) ([]byte, error) {
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

	content := getFile(baseDir, filename)

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

// getFile get file content with no error, return empty if something wrong
//
//	bytes, err := file.ReadDirect("mock/mock.json")
//
func getFile(baseDir, filename string) []byte {
	fullPath := Lookup(baseDir, filename)
	if fullPath == "" {
		return nil
	}
	content, err := ReadDirect(fullPath)
	if err != nil {
		return nil
	}
	return content
}

// Read binary data direct from file
//
//	bytes, err := file.ReadDirect("mock/mock.json")
//
func ReadDirect(fullPath string) ([]byte, error) {
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
	return content, nil
}

// ReadText read text from file, filename can be relative path
//
//	txt, err := file.ReadText("mock/mock.json", NoCache, 0)
//
func ReadText(baseDir, filename string, cacheType CacheType, d time.Duration) (string, error) {
	bytes, err := Read(baseDir, filename, cacheType, d)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ReadJSON read json object from file, filename can be relative path, return nil if file not found
//
//	f, err := file.ReadJSON("mock/mock.json", NoCache, 0)
//
func ReadJSON(baseDir, filename string, cacheType CacheType, d time.Duration) (map[string]interface{}, error) {
	bytes, err := Read(baseDir, filename, cacheType, d)
	if err != nil {
		return nil, err
	}
	content := make(map[string]interface{})
	if bytes == nil {
		return nil, nil
	}
	if err := json.Unmarshal([]byte(bytes), &content); err != nil {
		return nil, errors.Wrapf(err, "decode json %v", filename)
	}
	return content, nil
}

// ReadJSONDirect read json direct from file, filename can be relative path
//
//	f, err := file.ReadJSONDirect("mock/mock.json")
//
func ReadJSONDirect(filename string) (map[string]interface{}, error) {
	bytes, err := ReadDirect(filename)
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
//	found := Exists("mock/mock.json")
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
func Lookup(baseDir, name string) string {
	switch baseDir {
	case KeysDir:
		if keysBase != "" {
			return path.Join(keysBase, path.Join(baseDir, name))
		}
	case AssetsDir:
		if assetsBase != "" {
			return path.Join(assetsBase, path.Join(baseDir, name))
		}
	}

	curdir, err := os.Getwd()
	if err != nil {
		return ""
	}
	name = path.Join(baseDir, name)
	var filepath string
	dir := curdir
	for i := 0; i <= 5; i++ {
		filepath = path.Join(dir, name)
		if Exists(filepath) {
			// target found. we know the base now
			switch baseDir {
			case KeysDir:
				keysBase = dir
				//				fmt.Printf("keysBaseDir dir is %v \n", dir)
			case AssetsDir:
				assetsBase = dir
				//				fmt.Printf("assetsBaseDir dir is %v \n", dir)
			}
			return filepath
		}

		// already root dir, just give up
		if dir == "/" {
			break
		}
		//dir not exist, go up
		dir = path.Join(dir, "../")
	}
	return ""
}

// Key return bytes from keys/file, no cache for key file cause it always have another cache
//
//	bytes, err := Key("gcloud.json")
//
func Key(filename string) ([]byte, error) {
	return Read(KeysDir, filename, NoCache, 0)
}

// KeyJSON read json object from file, no cache for key file cause it always have another cache
//
//	f, err := KeyJSON("mock/mock.json")
//
func KeyJSON(filename string) (map[string]interface{}, error) {
	return ReadJSON(KeysDir, filename, NoCache, 0)
}

// KeyText read json object from file, no cache for key file cause it always have another cache
//
//	f, err := KeyText("mock/mock.json")
//
func KeyText(filename string) (string, error) {
	return ReadText(KeysDir, filename, NoCache, 0)
}

// I18nJSON return i18n file, use gzip cache and default cache time
//
//	j, err := I18nJSON("mock.json",0)
//
func I18nJSON(filename string, d time.Duration) (map[string]interface{}, error) {
	fullPath := path.Join(I18nDir, filename)
	return ReadJSON(AssetsDir, fullPath, GzipCache, d)
}

// I18nText return i18n text content, use gzip cache and default cache time
//
//	j, err := I18nText("mock.json",0)
//
func I18nText(filename string, d time.Duration) (string, error) {
	fullPath := path.Join(I18nDir, filename)
	return ReadText(AssetsDir, fullPath, GzipCache, d)
}
