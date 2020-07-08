package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestToInt(t *testing.T) {
	Convey("should convert interface{} to int", t, func() {
		var i interface{}
		i = 2
		num, err := ToInt(i)
		So(err, ShouldBeNil)
		So(num, ShouldEqual, 2)
	})
}

func TestToFloat64(t *testing.T) {
	Convey("should convert interface{} to float64", t, func() {
		var i interface{}
		i = 2
		num, err := ToFloat64(i)
		So(err, ShouldBeNil)
		So(num, ShouldEqual, 2)
	})
}

func TestToUint32(t *testing.T) {
	Convey("should convert interface{} to uint32", t, func() {
		var i interface{}
		i = 2
		num, err := ToUint32(i)
		So(err, ShouldBeNil)
		So(num, ShouldEqual, 2)
	})
}
