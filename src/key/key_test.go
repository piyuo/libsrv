package key

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldGetKeyContent(t *testing.T) {
	assert := assert.New(t)
	text, err := Text("gcloud.json")
	assert.Nil(err)
	assert.NotEmpty(text)

	bytes, err := Bytes("gcloud.json")
	assert.Nil(err)
	assert.NotNil(bytes)

	json, err := JSON("gcloud.json")
	assert.Nil(err)
	assert.NotEmpty(json["project_id"])
}

func TestReturnErrorWhenKeyNotExists(t *testing.T) {
	assert := assert.New(t)
	content, err := Text("not exists")
	assert.NotNil(err)
	assert.Empty(content)

	json, err := JSON("not exists")
	assert.NotNil(err)
	assert.Nil(json)

	bytes, err := Bytes("not exists")
	assert.NotNil(err)
	assert.Nil(bytes)

}

func TestConcurrentKey(t *testing.T) {
	var concurrent = 10
	var wg sync.WaitGroup
	wg.Add(concurrent)
	encryptDecrypt := func() {
		for i := 0; i < 10; i++ {
			_, err := Text("gcloud.json")
			if err != nil {
				t.Errorf("err should be nil, got %v", err)
			}
			_, err = JSON("gcloud.json")
			if err != nil {
				t.Errorf("err should be nil, got %v", err)
			}
			_, err = Bytes("gcloud.json")
			if err != nil {
				t.Errorf("err should be nil, got %v", err)
			}
			//fmt.Print(text + "\n")
		}
		wg.Done()
	}

	//create go routing to do counting
	for i := 0; i < concurrent; i++ {
		go encryptDecrypt()
	}
	wg.Wait()
}
