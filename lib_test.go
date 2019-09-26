package libsrv

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLib(t *testing.T) {
	text := lib()
	Convey("lib should return hello world", t, func() {
		So(text, ShouldEqual, "hello world")
	})
}
