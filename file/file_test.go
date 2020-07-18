package file

import (
	"os"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFile(t *testing.T) {
	currentDir, _ := os.Getwd()
	keyPath := path.Join(currentDir, "../../assets/key/gcloud.json")
	f, _ := Open(keyPath)
	defer f.Close()

	Convey("should have text'", t, func() {
		text := f.Text()
		So(len(text), ShouldBeGreaterThan, 1)
	})

	Convey("should have json'", t, func() {
		json, err := f.JSON()
		So(err, ShouldBeNil)
		So(json["project_id"], ShouldEqual, "piyuo-beta")
	})

	Convey("should have bytes'", t, func() {
		bytes := f.Bytes()
		So(len(bytes), ShouldBeGreaterThan, 0)
	})

}
