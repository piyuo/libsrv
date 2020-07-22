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

	Convey("should have text", t, func() {
		text, err := ReadText(keyPath)
		So(err, ShouldBeNil)
		So(len(text), ShouldBeGreaterThan, 1)
	})

	Convey("should have json", t, func() {
		json, err := ReadJSON(keyPath)
		So(err, ShouldBeNil)
		So(json["project_id"], ShouldEqual, "piyuo-beta")
	})
}

func TestFindDir(t *testing.T) {
	Convey("should find assets", t, func() {
		dir, found := FindDir("assets")
		So(found, ShouldBeTrue)
		So(dir, ShouldNotBeEmpty)
		dir, found = FindDir("not exist")
		So(found, ShouldBeFalse)
		So(dir, ShouldBeEmpty)
	})

}
