package crypto

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrypto(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	crypted, err := Encrypt("hello")
	assert.Nil(err)

	crypted1, err := Encrypt("hello1")
	assert.Nil(err)

	assert.NotEmpty(crypted)
	assert.NotEmpty(crypted1)
	result, err := Decrypt(crypted)
	assert.Nil(err)

	result1, err := Decrypt(crypted1)
	assert.Nil(err)

	assert.Equal("hello", result)
	assert.Equal("hello1", result1)
}

func TestShouldReturnErrorDecryptWrongString(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	_, err := Decrypt("")
	assert.NotNil(err)
	_, err1 := Decrypt("something wrong")
	assert.NotNil(err1)
}

func TestConcurrentCrypto(t *testing.T) {
	t.Parallel()
	var concurrent = 20
	var wg sync.WaitGroup
	wg.Add(concurrent)
	encryptDecrypt := func() {
		for i := 0; i < 100; i++ {
			id := strconv.Itoa(rand.Intn(1000))
			crypted, err := Encrypt("hello" + id)
			if err != nil {
				t.Errorf("err should be nil, got %v", err)
			}
			//fmt.Print(crypted + "\n")

			result, err := Decrypt(crypted)
			if err != nil {
				t.Errorf("err should be nil, got %v", err)
			}
			if result != "hello"+id {
				t.Errorf("failed to decrypt result, got %v", result)
			}
			//fmt.Print(result + "\n")

		}
		wg.Done()
	}

	//create go routing to do counting
	for i := 0; i < concurrent; i++ {
		go encryptDecrypt()
	}
	wg.Wait()
}
