package key

import (
	"errors"
	"path"

	"github.com/piyuo/libsrv/cache"
	"github.com/piyuo/libsrv/file"
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
	cache.Set(cache.HIGH, cachename, text)
	return text, nil
}

// TextWithoutCache return key text content, no cache on return value
//	text, err := key.TextWithoutCache("log.json")
//
func TextWithoutCache(name string) (string, error) {
	keyFile := path.Join("keys", name)
	keypath, found := file.Find(keyFile)
	if !found {
		return "", errors.New(keyFile + " not found")
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

	cache.Set(cache.HIGH, cachename, json)
	return json, nil
}

// JSONWithoutCache return key json object, no cache on return value
//
//	json, err := key.JSONWithoutCache("log.json")
//
func JSONWithoutCache(name string) (map[string]interface{}, error) {
	keyFile := path.Join("keys", name)
	keypath, found := file.Find(keyFile)
	if !found {
		return nil, errors.New(keyFile + " not found")
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

	cache.Set(cache.HIGH, cachename, bytes)
	return bytes, nil
}

// BytesWithoutCache return key bytes from /keys, no cache on return value
//
//	bytes, err := key.BytesWithoutCache("log.json")
//
func BytesWithoutCache(name string) ([]byte, error) {
	keyFile := path.Join("keys", name)
	keypath, found := file.Find(keyFile)
	if !found {
		return nil, errors.New(keyFile + " not found")
	}

	bytes, err := file.Read(keypath)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
