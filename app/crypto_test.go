package app

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCrypto(t *testing.T) {
	crypto := NewCrypto()
	Convey("should encrypt decrypt string", t, func() {
		crypted, _ := crypto.Encrypt("hello")
		crypted1, _ := crypto.Encrypt("hello1")
		So(crypted, ShouldNotBeEmpty)
		So(crypted1, ShouldNotBeEmpty)
		result, _ := crypto.Decrypt(crypted)
		result1, _ := crypto.Decrypt(crypted1)
		So(result, ShouldEqual, "hello")
		So(result1, ShouldEqual, "hello1")
	})

	Convey("should has error when decrypt empty or wrong string", t, func() {
		_, err := crypto.Decrypt("")
		So(err, ShouldNotBeNil)
		_, err1 := crypto.Decrypt("something wrong")
		So(err1, ShouldNotBeNil)
	})
}
