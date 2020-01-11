package libsrv

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJSONFIle(t *testing.T) {

	jsonfile, err := NewJSONFile(Sys().JoinCurrentDir("keys/log.key"))
	if err != nil {
		panic(err)
	}
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
