package file

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
)

// Read binary data from file, filename can be relative path
//
//	bytes, err := file.Read("mock/mock.json")
//	if( err != nil){
//		return err
//	}
//	So(len(bytes), ShouldBeGreaterThan, 0)
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
//	if( err != nil){
//		return err
//	}
//	So(txt, ShouldEqual, "")
//
func ReadText(filename string) (string, error) {
	bytes, err := Read(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ReadJSON read json object from file, filename can be relative path
//
//	f, err := file.ReadJSON("mock/mock.json")
//	if( err != nil){
//		return err
//	}
//	So(json["project_id"], ShouldEqual, "piyuo-beta")
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

// FindDir find dir from current path all the way to the top, return actual path where dir locate
//
//	path, err := FindDir("keys")
//
func FindDir(dirname string) (string, bool) {
	curdir, err := os.Getwd()
	if err != nil {
		return "", false
	}

	var dir string
	for i := 0; i <= 5; i++ {
		dir = path.Join(curdir, dirname)
		if _, err = os.Stat(dir); err == nil {
			//dir exist
			return dir, true
		}
		//dir not exist, go up
		curdir = path.Join(curdir, "../")
	}
	return "", false
}
