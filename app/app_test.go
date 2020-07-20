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
		tenMinutesLater := time.Now().Add(24 * time.Hour)
		So(dateline.Before(tenMinutesLater), ShouldBeTrue)
	})
}

func TestIsSlow(t *testing.T) {
	Convey("should determine slow work'", t, func() {
		// 3 seconds execution time is not slow
		So(IsSlow(5000), ShouldEqual, 0)
		// 20 seconds execution time is really slow
		So(IsSlow(20000000), ShouldBeGreaterThan, 5000)
	})
}

func TestBasicFunction(t *testing.T) {
	Convey("should able join dir and current dir'", t, func() {
		text, err := JoinCurrentDir("../../../")
		So(err, ShouldBeNil)
		So(strings.HasSuffix(text, "/go"), ShouldBeTrue)
	})
	Convey("should set env PIYUO_APP'", t, func() {
		id := PiyuoID()
		So(id, ShouldNotBeEmpty)
	})
	Convey("should get key path'", t, func() {
		path, err := KeyPath("gcloud")
		So(err, ShouldBeNil)
		So(strings.HasSuffix(path, "/gcloud.json"), ShouldBeTrue)
	})
	Convey("should get region key path'", t, func() {
		path, err := RegionKeyPath("us")
		So(err, ShouldBeNil)
		So(strings.HasSuffix(path, "/region/us.json"), ShouldBeTrue)
	})
	Convey("should get key content'", t, func() {
		text, err := Key("gcloud")
		So(err, ShouldBeNil)
		So(text, ShouldNotBeEmpty)
	})
	Convey("should get region key content'", t, func() {
		text, err := RegionKey("us")
		So(err, ShouldBeNil)
		So(text, ShouldNotBeEmpty)
	})
}
