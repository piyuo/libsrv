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
		So(s, ShouldNotBeEmpty)

		m2 := MapFromString(s)
		So(m2["a"], ShouldEqual, "1")
		So(m2["b"], ShouldEqual, "2")
	})
}

func TestMapEmpty(t *testing.T) {
	Convey("should allow empty map", t, func() {
		m := map[string]string{}
		s := MapToString(m)
		So(s, ShouldBeEmpty)

		m2 := MapFromString(s)
		So(len(m2), ShouldEqual, 0)
	})
}

func TestMapEmpty2(t *testing.T) {
	Convey("should allow empty map", t, func() {
		m := MapFromString("=")
		So(len(m), ShouldEqual, 0)
	})
}
