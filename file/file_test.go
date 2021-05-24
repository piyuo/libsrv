package file

import (
	"testing"

	"github.com/piyuo/libsrv/google"
	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	keyFile := "gcloud-test.json"
	//should have bytes
	bytes, err := Read(KeysDir, keyFile, NoCache, 0)
	assert.Nil(err)
	assert.NotEmpty(bytes)
	//should not have bytes
	bytes, err = Read(KeysDir, "not exist2", NoCache, 0)
	assert.NotNil(err)
	assert.Nil(bytes)
	//should have text
	text, err := ReadText(KeysDir, keyFile, NoCache, 0)
	assert.Nil(err)
	assert.NotEmpty(text, 1)
	//should not have bytes
	text, err = ReadText(KeysDir, "not exist", NoCache, 0)
	assert.NotNil(err)
	assert.Empty(text)
	//should have json
	json, err := ReadJSON(KeysDir, keyFile, NoCache, 0)
	assert.Nil(err)
	assert.Equal(google.TestProject, json["project_id"])
	//should not have json
	json, err = ReadJSON(KeysDir, "not exist", NoCache, 0)
	assert.NotNil(err)
	assert.Nil(json)
}

func TestDirect(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	keyPath := Lookup(KeysDir, "gcloud-test.json")

	bytes, err := ReadDirect(keyPath)
	assert.Nil(err)
	assert.NotEmpty(bytes)

	json, err := ReadJSONDirect(keyPath)
	assert.Nil(err)
	assert.NotNil(json)
}

func TestCache(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	keyFile := "gcloud-test.json"
	bytes, err := Read(KeysDir, keyFile, Cache, 0)
	assert.Nil(err)
	bytes2, err := Read(KeysDir, keyFile, Cache, 0)
	assert.Nil(err)
	assert.Equal(&bytes, &bytes2)

	text, err := ReadText(KeysDir, keyFile, Cache, 0)
	assert.Nil(err)
	assert.NotEmpty(text)

	json, err := ReadJSON(KeysDir, keyFile, Cache, 0)
	assert.Nil(err)
	assert.NotEmpty(json)
}

func TestGzipCache(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	keyFile := "gcloud.json"
	bytes, err := Read(KeysDir, keyFile, GzipCache, 0)
	assert.Nil(err)
	bytes2, err := Read(KeysDir, keyFile, GzipCache, 0)
	assert.Nil(err)
	assert.Equal(bytes, bytes2)

	text, err := ReadText(KeysDir, keyFile, GzipCache, 0)
	assert.Nil(err)
	assert.NotEmpty(text)

	json, err := ReadJSON(KeysDir, keyFile, GzipCache, 0)
	assert.Nil(err)
	assert.NotEmpty(json)
}

func TestLookup(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	dir := Lookup(KeysDir, "not-exist") // can't determine
	assert.Empty(dir)
	assert.False(Exists(dir))

	dir = Lookup(KeysDir, "gcloud-test.json")
	assert.NotEmpty(dir)
	assert.True(Exists(dir))

	dir = Lookup(KeysDir, "not-exist") // have base dir can determine
	assert.NotEmpty(dir)
	assert.False(Exists(dir))
}

func TestKey(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	bytes, err := Key("gcloud.json")
	assert.Nil(err)
	assert.NotEmpty(bytes)
	j, err := KeyJSON("gcloud.json")
	assert.Nil(err)
	assert.NotNil(j)
	txt, err := KeyText("gcloud.json")
	assert.Nil(err)
	assert.NotEmpty(txt)
}

func TestI18n(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	bytes, err := I18nText("mock_en_US.json", 0)
	assert.Nil(err)
	assert.NotEmpty(bytes)
	j, err := I18nJSON("mock_en_US.json", 0)
	assert.Nil(err)
	assert.NotNil(j)
}
