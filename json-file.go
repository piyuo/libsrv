package libsrv

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// JSONFile represent  json file
type JSONFile interface {
	Load(filename string) error
	JSON() (map[string]interface{}, error)
	Text() (string, error)
	Close()
}

//NewJSONFile create JSONFile instance
//
//	jsonfile, err := NewJSONFile("mock/mock.key")
//	if( err != nil){
//		panic(err)
//	}
//	defer jsonfile.Close()
func NewJSONFile(filename string) (JSONFile, error) {
	jf := &jsonfile{}
	err := jf.Load(filename)
	if err != nil {
		return nil, err
	}
	return jf, nil
}

type jsonfile struct {
	jsonFile *os.File
	json     map[string]interface{}
	text     string
	bytes    []byte
}

func (j *jsonfile) Load(filename string) error {
	file, err := os.Open(filename)

	if err != nil {
		return errors.Wrap(err, filename+" can not open")
	}
	defer file.Close()

	j.bytes, err = ioutil.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, filename+" can not read")
	}
	return nil
}

// JSON get json data store
//
//	fmt.Println(JSON()["users"])
func (j *jsonfile) JSON() (map[string]interface{}, error) {
	if j.json != nil {
		return j.json, nil
	}

	if j.bytes == nil {
		return nil, errors.New("JSON() need Load() first")
	}

	err := json.Unmarshal([]byte(j.bytes), &j.json)
	if err != nil {
		return nil, errors.Wrap(err, "decode json fail")
	}
	return j.json, nil
}

// Text get json data store
//
//	fmt.Println(Text())
func (j *jsonfile) Text() (string, error) {
	if j.text != "" {
		return j.text, nil
	}

	if j.bytes == nil {
		return "", errors.New("Text() need Load() first")
	}

	j.text = string(j.bytes)
	return j.text, nil
}

// Close release memory
//
//	fmt.Close()
func (j *jsonfile) Close() {
	j.bytes = nil
	j.json = nil
	j.text = ""
}
