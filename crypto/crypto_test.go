package crypto

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCrypto(t *testing.T) {
	Convey("should encrypt decrypt string", t, func() {
		crypted, err := Encrypt("hello")
		So(err, ShouldBeNil)

		crypted1, err := Encrypt("hello1")
		So(err, ShouldBeNil)

		So(crypted, ShouldNotBeEmpty)
		So(crypted1, ShouldNotBeEmpty)
		result, err := Decrypt(crypted)
		So(err, ShouldBeNil)

		result1, err := Decrypt(crypted1)
		So(err, ShouldBeNil)

		So(result, ShouldEqual, "hello")
		So(result1, ShouldEqual, "hello1")
	})

	Convey("should has error when decrypt empty or wrong string", t, func() {
		_, err := Decrypt("")
		So(err, ShouldNotBeNil)
		_, err1 := Decrypt("something wrong")
		So(err1, ShouldNotBeNil)
	})
}

func TestConcurrentCrypto(t *testing.T) {
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
