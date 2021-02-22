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
		return nil, errors.Wrap(err, "failed to open file: "+filename)
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
		return nil, errors.Wrap(err, "failed to decode json: "+filename)
	}
	return content, nil
}

// Find find dir or file from current path all the way to the top, return actual path where dir or file locate
//
//	dir, err := Find("keys")
//
func Find(name string) (string, bool) {
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

	filepath = path.Join(curdir, name)
	fmt.Printf("failed to find %v in %v or parent dir\n", name, curdir)
	return "", false
}
