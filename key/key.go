package key

import (
	"os"
	"path"

	cache "github.com/piyuo/libsrv/cache"
	file "github.com/piyuo/libsrv/file"
	"github.com/pkg/errors"
)

// getPath get key real path from name, key path is "keys/" which can be place under /src/keys or /src/project/keys
//
//	path, err := getPath("log.json")
//
func getPath(keyname string) (string, error) {
	curdir, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "failed to Getwd()")
	}

	var keydir string
	var keypath string
	for i := 0; i <= 3; i++ {
		keydir = path.Join(curdir, "keys")
		keypath = path.Join(keydir, keyname)
		if _, err = os.Stat(keypath); err == nil {
			//keyPath exist
			return keypath, nil
		}
		//keyPath not exist, go up
		curdir = path.Join(curdir, "../")
	}
	return "", errors.New("failed to find " + keyname + ".json in keys/ or ../keys/")
}

// Text return key text content, return key content wil be cache to reuse in the future
//
//	key, err := key.Text("log.json")
//
func Text(name string) (string, error) {
	keyname := name + "Text"
	value, found := cache.Get(keyname)
	if found {
		return value.(string), nil
	}

	keyPath, err := getPath(name)
	if err != nil {
		return "", err
	}
	text, err := file.ReadText(keyPath)
	if err != nil {
		return "", err
	}
	cache.Set(keyname, text, -1) // key never expire, cause we always need it
	return text, nil
}

// JSON return key json object, return key content wil be cache to reuse in the future
//
//	key, err := keys.JSON("log.json")
//
func JSON(name string) (map[string]interface{}, error) {
	keyname := name + "JSON"
	value, found := cache.Get(keyname)
	if found {
		return value.(map[string]interface{}), nil
	}

	keyPath, err := getPath(name)
	if err != nil {
		return nil, err
	}
	json, err := file.ReadJSON(keyPath)
	if err != nil {
		return nil, err
	}
	cache.Set(keyname, json, -1) // key never expire, cause we always need it
	return json, nil
}

// Bytes return key bytes from /keys, return key content wil be cache to reuse in the future
//	key, err := keys.Key("log.json")
//
func Bytes(name string) ([]byte, error) {
	keyname := name + "Byte"
	value, found := cache.Get(keyname)
	if found {
		return value.([]byte), nil
	}

	keyPath, err := getPath(name)
	if err != nil {
		return nil, err
	}
	bytes, err := file.Read(keyPath)
	if err != nil {
		return nil, err
	}
	cache.Set(keyname, bytes, -1) // key never expire, cause we always need it
	return bytes, nil
}
