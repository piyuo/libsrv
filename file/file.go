package file

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// File represent  json file
type File interface {
	JSON() (map[string]interface{}, error)
	Text() string
	Bytes() []byte
	Close()
}

//Open file to read JSON or Text
//
//	f, err := file.Open("mock/mock.key")
//	if( err != nil){
//		panic(err)
//	}
//	defer f.Close()
func Open(filename string) (File, error) {
	f := &file{}
	osFile, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open "+filename)
	}
	defer osFile.Close()

	f.bytes, err = ioutil.ReadAll(osFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read "+filename)
	}
	return f, nil
}

type file struct {
	f     *os.File
	json  map[string]interface{}
	text  string
	bytes []byte
}

// JSON get json data store
//
//	fmt.Println(f.JSON()["users"])
func (f *file) JSON() (map[string]interface{}, error) {
	if f.json != nil {
		return f.json, nil
	}

	err := json.Unmarshal([]byte(f.bytes), &f.json)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode json")
	}
	return f.json, nil
}

// Text get json data store
//
//	f.Text()
func (f *file) Text() string {
	if f.text != "" {
		return f.text
	}
	f.text = string(f.bytes)
	return f.text
}

//Bytes return bytes
//
//	f.Bytes()
func (f *file) Bytes() []byte {
	return f.bytes
}

// Close release memory
//
//	f.Close()
func (f *file) Close() {
	f.bytes = nil
	f.json = nil
	f.text = ""
}
