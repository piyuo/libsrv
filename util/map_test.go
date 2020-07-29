package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMap(t *testing.T) {
	Convey("should ToString() and FromString", t, func() {
		m := map[string]string{
			"a": "1",
			"b": "2",
		}
		s := MapToString(m)
		So(s, ShouldEqual, "a=1&b=2")

		m2 := MapFromString(s)
		So(m2["a"], ShouldEqual, "1")
		So(m2["b"], ShouldEqual, "2")
	})
}
