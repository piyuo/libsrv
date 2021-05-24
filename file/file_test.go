package file

import (
	"testing"

	"github.com/piyuo/libsrv/google"
	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	keyPath := "keys/gcloud-test.json"
	//should have bytes
	bytes, err := Read(keyPath, NoCache, 0)
	assert.Nil(err)
	assert.NotEmpty(bytes)
	//should not have bytes
	bytes, err = Read("not exist", NoCache, 0)
	assert.NotNil(err)
	assert.Nil(bytes)
	//should have text
	text, err := ReadText(keyPath, NoCache, 0)
	assert.Nil(err)
	assert.NotEmpty(text, 1)
	//should not have bytes
	text, err = ReadText("not exist", NoCache, 0)
	assert.NotNil(err)
	assert.Empty(text)
	//should have json
	json, err := ReadJSON(keyPath, NoCache, 0)
	assert.Nil(err)
	assert.Equal(google.TestProject, json["project_id"])
	//should not have json
	json, err = ReadJSON("not exist", NoCache, 0)
	assert.NotNil(err)
	assert.Nil(json)
}

func TestCache(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	keyPath := "keys/gcloud-test.json"
	bytes, err := Read(keyPath, Cache, 0)
	assert.Nil(err)
	bytes2, err := Read(keyPath, Cache, 0)
	assert.Nil(err)
	assert.Equal(&bytes, &bytes2)

	text, err := ReadText(keyPath, Cache, 0)
	assert.Nil(err)
	assert.NotEmpty(text)

	json, err := ReadJSON(keyPath, Cache, 0)
	assert.Nil(err)
	assert.NotEmpty(json)
}

func TestGzipCache(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	keyPath := "keys/gcloud.json"
	bytes, err := Read(keyPath, GzipCache, 0)
	assert.Nil(err)
	bytes2, err := Read(keyPath, GzipCache, 0)
	assert.Nil(err)
	assert.Equal(bytes, bytes2)

	text, err := ReadText(keyPath, GzipCache, 0)
	assert.Nil(err)
	assert.NotEmpty(text)

	json, err := ReadJSON(keyPath, GzipCache, 0)
	assert.Nil(err)
	assert.NotEmpty(json)
}

func TestLookup(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	dir := Lookup("not-exist") // can't determine
	assert.Empty(dir)
	assert.False(Exists(dir))

	dir = Lookup("keys/gcloud-test.json")
	assert.NotEmpty(dir)
	assert.True(Exists(dir))

	dir = Lookup("not-exist") // have base dir can determine
	assert.NotEmpty(dir)
	assert.False(Exists(dir))
}
