package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
)

// Read binary data from file, filename can be relative path
//
//	bytes, err := file.Read("mock/mock.json")
//
func Read(filename string) ([]byte, error) {
	osFile, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "open file: %v", filename)
	}
	defer osFile.Close()
	return ioutil.ReadAll(osFile)
}

// ReadText read text from file, filename can be relative path
//
//	txt, err := file.ReadText("mock/mock.json")
//
func ReadText(filename string) (string, error) {
	bytes, err := Read(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// WriteText write text to file, filename can be relative path
//
//	txt, err := file.WriteText("hello.txt","hello")
//
func WriteText(filename, text string) error {
	bytes := []byte(text)
	return ioutil.WriteFile(filename, bytes, 0644)
}

// ReadJSON read json object from file, filename can be relative path
//
//	f, err := file.ReadJSON("mock/mock.json")
//
func ReadJSON(filename string) (map[string]interface{}, error) {
	bytes, err := Read(filename)
	if err != nil {
		return nil, err
	}
	content := make(map[string]interface{})
	if err := json.Unmarshal([]byte(bytes), &content); err != nil {
		return nil, errors.Wrapf(err, "decode json %v", filename)
	}
	return content, nil
}

// Lookup find dir or file from current path all the way to the top, return actual path where dir or file locate
//
//	dir, err := Lookup("keys")
//
func Lookup(name string) (string, bool) {
	curdir, err := os.Getwd()
	if err != nil {
		return "", false
	}

	var filepath string
	dir := curdir
	for i := 0; i <= 5; i++ {
		filepath = path.Join(dir, name)
		if _, err = os.Stat(filepath); err == nil {
			//dir exist
			return filepath, true
		}

		//root dir, just give up
		if dir == "/" {
			break
		}

		//dir not exist, go up
		dir = path.Join(dir, "../")
	}
	fmt.Printf("%v not found in %v or parent dir", name, curdir)
	return "", false
}
