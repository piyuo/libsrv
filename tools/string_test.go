package libsrv

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestString(t *testing.T) {
	Convey("should work on string'", t, func() {
		So(StringBetween("123", "1", "3"), ShouldEqual, "2")
		So(StringBefore("123", "2"), ShouldEqual, "1")
		So(StringAfter("123", "2"), ShouldEqual, "3")
	})
}
