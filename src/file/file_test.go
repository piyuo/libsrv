package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	assert := assert.New(t)
	keyPath := "../../../keys/gcloud.json"
	//should have bytes
	bytes, err := Read("../../../keys/gcloud.json")
	assert.Nil(err)
	assert.Greater(len(bytes), 0)
	//should not have bytes
	bytes, err = Read("not exist")
	assert.NotNil(err)
	assert.Nil(bytes)
	//should have text
	text, err := ReadText(keyPath)
	assert.Nil(err)
	assert.Greater(len(text), 1)
	//should not have bytes
	text, err = ReadText("not exist")
	assert.NotNil(err)
	assert.Empty(text)
	//should have json
	json, err := ReadJSON(keyPath)
	assert.Nil(err)
	assert.Equal("piyuo-beta", json["project_id"])
	//should not have json
	json, err = ReadJSON("not exist")
	assert.NotNil(err)
	assert.Nil(json)
}

func TestFileReadWrite(t *testing.T) {
	assert := assert.New(t)
	filename := "test.txt"
	err := WriteText(filename, "hello")
	assert.Nil(err)
	text, err := ReadText(filename)
	assert.Nil(err)
	assert.Equal("hello", text)
	defer os.Remove(filename)
}

func TestFileFind(t *testing.T) {
	assert := assert.New(t)
	dir, found := Find("assets")
	assert.True(found)
	assert.NotEmpty(dir)
	dir, found = Find("not exist")
	assert.False(found)
	assert.Empty(dir)
}
