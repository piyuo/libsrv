package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestString(t *testing.T) {
	Convey("should work on string", t, func() {
		So(StringBetween("123", "1", "3"), ShouldEqual, "2")
		So(StringBetween("123", "a", "3"), ShouldEqual, "")
		So(StringBetween("123", "1", "a"), ShouldEqual, "")
		So(StringBetween("111", "1", "1"), ShouldEqual, "")

		So(StringBefore("123", "2"), ShouldEqual, "1")
		So(StringBefore("123", "a"), ShouldEqual, "")
		So(StringAfter("123", "2"), ShouldEqual, "3")
		So(StringAfter("123", "a"), ShouldEqual, "")
		So(StringAfter("111", "1"), ShouldEqual, "")
	})
}
