package app

import (
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCheck(t *testing.T) {
	Convey("should pass check with no panic", t, func() {
		Check()
		So(IsProduction(), ShouldBeFalse)
	})
	Convey("should return piyuo id'", t, func() {
		piyuoID := PiyuoID()
		So(piyuoID, ShouldNotBeEmpty)
	})
}

func TestDateline(t *testing.T) {
	Convey("should return deadline'", t, func() {
		dateline := ContextDateline()
		So(dateline.After(time.Now()), ShouldBeTrue)

		//dateline should not greater than 10 min.
		tenMinutesLater := time.Now().Add(10 * time.Minute)
		So(dateline.Before(tenMinutesLater), ShouldBeTrue)
	})
}

func TestIsSlow(t *testing.T) {
	Convey("should determine slow work'", t, func() {
		// 3 seconds execution time is not slow
		So(IsSlow(5000), ShouldEqual, 0)
		// 20 seconds execution time is really slow
		So(IsSlow(20000), ShouldEqual, 10000)
	})
}

func TestBasicFunction(t *testing.T) {
	Convey("should able join dir and current dir'", t, func() {
		text := JoinCurrentDir("../../../")
		So(strings.HasSuffix(text, "/go"), ShouldBeTrue)
	})
	Convey("should set env PIYUO_APP'", t, func() {
		id := PiyuoID()
		So(id, ShouldNotBeEmpty)
		So(IsProduction(), ShouldEqual, false)
	})
	Convey("should get key path'", t, func() {
		path, err := KeyPath("log-gcp")
		So(err, ShouldBeNil)
		So(strings.HasSuffix(path, "/log-gcp.key"), ShouldBeTrue)
	})
	Convey("should get key content'", t, func() {
		text, err := Key("log-gcp")
		So(err, ShouldBeNil)
		So(text, ShouldNotBeEmpty)
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
