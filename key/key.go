package key

import (
	"os"
	"path"

	cache "github.com/piyuo/libsrv/cache"
	file "github.com/piyuo/libsrv/file"
	"github.com/pkg/errors"
)

func getKeyPath(keyname string) (string, error) {
	keydir, found := file.FindDir("keys")
	if found {
		keyfile := path.Join(keydir, keyname)
		if _, err := os.Stat(keyfile); err == nil {
			//keyfile exist
			return keyfile, nil
		}
		return "", errors.New(keyname + " not found in /keys")
	}
	return "", errors.New("/keys dir not found")
}

// Text return key text content, return key content wil be cache to reuse in the future
//
//	key, err := key.Text("log.json")
//
func Text(name string) (string, error) {
	cachename := "KEY" + name + "TEXT"
	value, found := cache.Get(cachename)
	if found {
		return value.(string), nil
	}

	keypath, err := getKeyPath(name)
	if err != nil {
		return "", err
	}

	text, err := file.ReadText(keypath)
	if err != nil {
		return "", err
	}
	cache.Set(cachename, text, -1) // key never expire, cause we always need it
	return text, nil
}

// JSON return key json object, return key content wil be cache to reuse in the future
//
//	key, err := keys.JSON("log.json")
//
func JSON(name string) (map[string]interface{}, error) {
	cachename := "KEY" + name + "JSON"
	value, found := cache.Get(cachename)
	if found {
		return value.(map[string]interface{}), nil
	}

	keypath, err := getKeyPath(name)
	if err != nil {
		return nil, err
	}

	json, err := file.ReadJSON(keypath)
	if err != nil {
		return nil, err
	}
	cache.Set(cachename, json, -1) // key never expire, cause we always need it
	return json, nil
}

// Bytes return key bytes from /keys, return key content wil be cache to reuse in the future
//	key, err := keys.Key("log.json")
//
func Bytes(name string) ([]byte, error) {
	cachename := "KEY" + name + "BYTE"
	value, found := cache.Get(cachename)
	if found {
		return value.([]byte), nil
	}

	keypath, err := getKeyPath(name)
	if err != nil {
		return nil, err
	}

	bytes, err := file.Read(keypath)
	if err != nil {
		return nil, err
	}
	cache.Set(cachename, bytes, -1) // key never expire, cause we always need it
	return bytes, nil
}
