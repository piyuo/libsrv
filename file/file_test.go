package file

import (
	"os"
	"testing"

	"github.com/piyuo/libsrv/google"
	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	keyPath := "../../keys/gcloud-test.json"
	//should have bytes
	bytes, err := Read("../../keys/gcloud-test.json")
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
	assert.Equal(google.TestProject, json["project_id"])
	//should not have json
	json, err = ReadJSON("not exist")
	assert.NotNil(err)
	assert.Nil(json)
}

func TestReadWrite(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	filename := "test.txt"
	err := WriteText(filename, "hello")
	assert.Nil(err)
	text, err := ReadText(filename)
	assert.Nil(err)
	assert.Equal("hello", text)
	defer os.Remove(filename)
}

func TestLookup(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	dir, found := Lookup("assets")
	assert.True(found)
	assert.NotEmpty(dir)
	dir, found = Lookup("not exist")
	assert.False(found)
	assert.Empty(dir)
}
