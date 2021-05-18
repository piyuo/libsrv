package key

import (
	"encoding/json"
	"path"

	"github.com/pkg/errors"

	"github.com/piyuo/libsrv/cache"
	"github.com/piyuo/libsrv/file"
)

// Text return key text content, return key content wil be cache to reuse in the future
//
//	text, err := key.Text("log.json")
//
func Text(name string) (string, error) {
	cachename := "kT" + name
	found, value, err := cache.GetString(cachename)
	if err != nil {
		return "", errors.Wrap(err, "get cache "+cachename)
	}
	if found {
		return value, nil
	}

	text, err := TextWithoutCache(name)
	if err != nil {
		return "", errors.Wrap(err, "get text without cache")
	}
	if err := cache.SetString(cachename, text, 0); err != nil {
		return "", errors.Wrap(err, "set cache "+cachename)
	}
	return text, nil
}

// TextWithoutCache return key text content, no cache on return value
//
//	text, err := key.TextWithoutCache("log.json")
//
func TextWithoutCache(name string) (string, error) {
	keyFile := path.Join("keys", name)
	keypath, found := file.Lookup(keyFile)
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
	cachename := "kJ" + name
	found, bytes, err := cache.Get(cachename)
	if err != nil {
		return nil, errors.Wrap(err, "get cache "+cachename)
	}
	if found {
		j := make(map[string]interface{})
		if err := json.Unmarshal(bytes, &j); err != nil {
			return nil, errors.Wrapf(err, "decode cache json %v", cachename)
		}
		return j, nil
	}

	bytes, err = BytesWithoutCache(name)
	if err != nil {
		return nil, errors.Wrap(err, "get bytes without cache")
	}
	if err := cache.Set(cachename, bytes, 0); err != nil {
		return nil, errors.Wrap(err, "set cache "+cachename)
	}
	j := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &j); err != nil {
		return nil, errors.Wrapf(err, "decode cache json %v", cachename)
	}
	return j, nil
}

// Bytes return key bytes from /keys, return key content wil be cache to reuse in the future
//
//	bytes, err := key.Bytes("log.json")
//
func Bytes(name string) ([]byte, error) {
	cachename := "kB" + name
	found, value, err := cache.Get(cachename)
	if err != nil {
		return nil, errors.Wrap(err, "get cache "+cachename)
	}
	if found {
		return value, nil
	}

	bytes, err := BytesWithoutCache(name)
	if err != nil {
		return nil, errors.Wrap(err, "get bytes without cache")
	}

	if err := cache.Set(cachename, bytes, 0); err != nil {
		return nil, errors.Wrap(err, "set cache "+cachename)
	}
	return bytes, nil
}

// BytesWithoutCache return key bytes from /keys, no cache on return value
//
//	bytes, err := key.BytesWithoutCache("log.json")
//
func BytesWithoutCache(name string) ([]byte, error) {
	keyFile := path.Join("keys", name)
	keypath, found := file.Lookup(keyFile)
	if !found {
		return nil, errors.New(keyFile + " not found")
	}

	bytes, err := file.Read(keypath)
	if err != nil {
		return nil, errors.Wrap(err, "read bytes "+keypath)
	}
	return bytes, nil
}
