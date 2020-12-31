package command

import (
	"testing"

	"github.com/piyuo/libsrv/command/shared"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShared(t *testing.T) {
	Convey("should create text response", t, func() {
		text := String("hi").(*shared.PbString)
		So(text.Value, ShouldEqual, "hi")
	})

	Convey("should create number response", t, func() {
		num := Int(201).(*shared.PbInt)
		So(num.Value, ShouldEqual, 201)
	})

	Convey("should create bool response", t, func() {
		b := Bool(true).(*shared.PbBool)
		So(b.Value, ShouldEqual, true)
		b = Bool(false).(*shared.PbBool)
		So(b.Value, ShouldEqual, false)
	})

	Convey("should create error response", t, func() {
		err := Error("errCode").(*shared.PbError)
		So(err.Code, ShouldEqual, "errCode")
	})

	Convey("should be OK", t, func() {
		ok := OK()
		So(IsOK(ok), ShouldBeTrue)
	})

	Convey("should be INVALID error", t, func() {
		err := Error("INVALID")
		So(IsError(err, "INVALID"), ShouldBeTrue)
	})

	Convey("should not be INVALID error", t, func() {
		So(IsError(nil, "INVALID"), ShouldBeFalse)
		err := 3
		So(IsError(err, "INVALID"), ShouldBeFalse)
	})
}

func TestPbString(t *testing.T) {
	Convey("should return is PbString", t, func() {
		So(IsString(nil, ""), ShouldBeFalse)
		So(IsString(String("hi"), ""), ShouldBeFalse)
		So(IsString(String("hi"), "hi"), ShouldBeTrue)
	})
}

func TestPbInt(t *testing.T) {
	Convey("should return is PbInt", t, func() {
		So(IsInt(nil, 1), ShouldBeFalse)
		So(IsInt(Int(12), 42), ShouldBeFalse)
		So(IsInt(Int(42), 42), ShouldBeTrue)
	})
}

func TestPbBool(t *testing.T) {
	Convey("should return is PbBool", t, func() {
		So(IsBool(nil, false), ShouldBeFalse)
		So(IsBool(Bool(false), true), ShouldBeFalse)
		So(IsBool(Bool(true), true), ShouldBeTrue)
	})
}
