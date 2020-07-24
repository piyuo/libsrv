package key

import (
	"errors"

	cache "github.com/piyuo/libsrv/cache"
	file "github.com/piyuo/libsrv/file"
)

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

	keypath, found := file.Find("keys/" + name)
	if !found {
		return "", errors.New("keys/" + name + " not found")
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

	keypath, found := file.Find("keys/" + name)
	if !found {
		return nil, errors.New("keys/" + name + " not found")
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

	keypath, found := file.Find("keys/" + name)
	if !found {
		return nil, errors.New("keys/" + name + " not found")
	}

	bytes, err := file.Read(keypath)
	if err != nil {
		return nil, err
	}
	cache.Set(cachename, bytes, -1) // key never expire, cause we always need it
	return bytes, nil
}
