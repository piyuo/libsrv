package key

import (
	"os"
	"path"
	"sync"

	file "github.com/piyuo/libsrv/file"
	"github.com/pkg/errors"
)

var cachedText = sync.Map{}
var cachedJSON = sync.Map{}
var cachedBytes = sync.Map{}

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
	result, ok := cachedText.Load(name)
	if ok {
		return result.(string), nil
	}

	keyPath, err := getPath(name)
	if err != nil {
		return "", err
	}
	text, err := file.ReadText(keyPath)
	if err != nil {
		return "", err
	}
	cachedText.Store(name, text)
	return text, nil
}

// JSON return key json object, return key content wil be cache to reuse in the future
//
//	key, err := keys.JSON("log.json")
//
func JSON(name string) (map[string]interface{}, error) {
	result, ok := cachedJSON.Load(name)
	if ok {
		return result.(map[string]interface{}), nil
	}

	keyPath, err := getPath(name)
	if err != nil {
		return nil, err
	}
	json, err := file.ReadJSON(keyPath)
	if err != nil {
		return nil, err
	}
	cachedJSON.Store(name, json)
	return json, nil
}

// Bytes return key bytes from /keys, return key content wil be cache to reuse in the future
//	key, err := keys.Key("log.json")
//
func Bytes(name string) ([]byte, error) {
	result, ok := cachedBytes.Load(name)
	if ok {
		return result.([]byte), nil
	}

	keyPath, err := getPath(name)
	if err != nil {
		return nil, err
	}
	bytes, err := file.Read(keyPath)
	if err != nil {
		return nil, err
	}
	cachedBytes.Store(name, bytes)
	return bytes, nil
}
