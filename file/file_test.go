package file

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFile(t *testing.T) {
	keyPath := "../../keys/gcloud.json"

	Convey("should have bytes", t, func() {
		bytes, err := Read("../../keys/gcloud.json")
		So(err, ShouldBeNil)
		So(len(bytes), ShouldBeGreaterThan, 0)
	})

	Convey("should not have bytes", t, func() {
		bytes, err := Read("not exist")
		So(err, ShouldNotBeNil)
		So(bytes, ShouldBeNil)
	})

	Convey("should have text", t, func() {
		text, err := ReadText(keyPath)
		So(err, ShouldBeNil)
		So(len(text), ShouldBeGreaterThan, 1)
	})

	Convey("should not have bytes", t, func() {
		text, err := ReadText("not exist")
		So(err, ShouldNotBeNil)
		So(text, ShouldBeEmpty)
	})

	Convey("should have json", t, func() {
		json, err := ReadJSON(keyPath)
		So(err, ShouldBeNil)
		So(json["project_id"], ShouldEqual, "piyuo-beta")
	})

	Convey("should not have json", t, func() {
		json, err := ReadJSON("not exist")
		So(err, ShouldNotBeNil)
		So(json, ShouldBeNil)
	})

}

func TestFind(t *testing.T) {
	Convey("should find assets", t, func() {
		dir, found := Find("assets")
		So(found, ShouldBeTrue)
		So(dir, ShouldNotBeEmpty)
		dir, found = Find("not exist")
		So(found, ShouldBeFalse)
		So(dir, ShouldBeEmpty)
	})

}
