package key

import (
	"errors"

	cache "github.com/piyuo/libsrv/cache"
	file "github.com/piyuo/libsrv/file"
)

// Text return key text content, return key content wil be cache to reuse in the future
//
//	text, err := key.Text("log.json")
//
func Text(name string) (string, error) {
	cachename := "KEY" + name + "TEXT"
	value, found := cache.Get(cachename)
	if found {
		return value.(string), nil
	}

	text, err := TextWithoutCache(name)
	if err != nil {
		return "", err
	}
	cache.Set(cachename, text, -1) // key never expire, cause we always need it
	return text, nil
}

// TextWithoutCache return key text content, no cache on return value
//	text, err := key.TextWithoutCache("log.json")
//
func TextWithoutCache(name string) (string, error) {
	keypath, found := file.Find("keys/" + name)
	if !found {
		return "", errors.New("keys/" + name + " not found")
	}

	text, err := file.ReadText(keypath)
	if err != nil {
		return "", err
	}
	return text, nil
}

// JSON return key json object, return key content wil be cache to reuse in the future
//
//	json, err := key.JSON("log.json")
//
func JSON(name string) (map[string]interface{}, error) {
	cachename := "KEY" + name + "JSON"
	value, found := cache.Get(cachename)
	if found {
		return value.(map[string]interface{}), nil
	}

	json, err := JSONWithoutCache(name)
	if err != nil {
		return nil, err
	}

	cache.Set(cachename, json, -1) // key never expire, cause we always need it
	return json, nil
}

// JSONWithoutCache return key json object, no cache on return value
//
//	json, err := key.JSONWithoutCache("log.json")
//
func JSONWithoutCache(name string) (map[string]interface{}, error) {

	keypath, found := file.Find("keys/" + name)
	if !found {
		return nil, errors.New("keys/" + name + " not found")
	}

	json, err := file.ReadJSON(keypath)
	if err != nil {
		return nil, err
	}
	return json, nil
}

// Bytes return key bytes from /keys, return key content wil be cache to reuse in the future
//
//	bytes, err := key.Bytes("log.json")
//
func Bytes(name string) ([]byte, error) {
	cachename := "KEY" + name + "BYTE"
	value, found := cache.Get(cachename)
	if found {
		return value.([]byte), nil
	}

	bytes, err := BytesWithoutCache(name)
	if err != nil {
		return nil, err
	}

	cache.Set(cachename, bytes, -1) // key never expire, cause we always need it
	return bytes, nil
}

// BytesWithoutCache return key bytes from /keys, no cache on return value
//
//	bytes, err := key.BytesWithoutCache("log.json")
//
func BytesWithoutCache(name string) ([]byte, error) {

	keypath, found := file.Find("keys/" + name)
	if !found {
		return nil, errors.New("keys/" + name + " not found")
	}

	bytes, err := file.Read(keypath)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
