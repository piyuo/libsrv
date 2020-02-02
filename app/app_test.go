package app

import (
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCheck(t *testing.T) {
	Convey("should pass check with no panic'", t, func() {
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
		dateline, err := ContextDateline()
		So(err, ShouldBeNil)
		So(dateline.After(time.Now()), ShouldBeTrue)

		//dateline should not greater than 10 min.
		tenMinutesLater := time.Now().Add(10 * time.Minute)
		So(dateline.Before(tenMinutesLater), ShouldBeTrue)
	})
}

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
