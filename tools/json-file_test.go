package tools

import (
	"os"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJSONFIle(t *testing.T) {
	currentDir, _ := os.Getwd()
	keyPath := path.Join(currentDir, "../keys/log.key")
	jsonfile, _ := NewJSONFile(keyPath)
	defer jsonfile.Close()

	Convey("should have text'", t, func() {
		text, _ := jsonfile.Text()
		So(len(text), ShouldBeGreaterThan, 1)
	})

	Convey("should have json'", t, func() {
		json, _ := jsonfile.JSON()
		So(json["project_id"], ShouldEqual, "master-255220")
	})

}
