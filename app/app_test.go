package app

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBasicFunction(t *testing.T) {
	Convey("should able join dir and current dir'", t, func() {
		text := JoinCurrentDir("../../../")
		So(strings.HasSuffix(text, "/go"), ShouldBeTrue)
	})
	Convey("should env PIYUO_ENV'", t, func() {
		id := PiyuoID()
		So(id, ShouldEqual, "piyuo-dev")
		So(IsProduction(), ShouldEqual, false)
	})
	Convey("should get key path'", t, func() {
		path, err := KeyPath("log-gcp")
		So(err, ShouldBeNil)
		So(strings.HasSuffix(path, "/log-gcp.key"), ShouldBeTrue)
	})
}

func TestAppCrypto(t *testing.T) {
	Convey("should encrypt decrypt string", t, func() {
		crypted, _ := Encrypt("hi")
		So(crypted, ShouldNotBeEmpty)
		result, _ := Decrypt(crypted)
		So(result, ShouldEqual, "hi")
	})

}
