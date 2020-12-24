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

	Convey("should be INVALID_MAIL error", t, func() {
		err := Error("INVALID_MAIL")
		So(IsError(err, "INVALID_MAIL"), ShouldBeTrue)
	})

	Convey("should not be INVALID_MAIL error", t, func() {
		So(IsError(nil, "INVALID_MAIL"), ShouldBeFalse)
		err := 3
		So(IsError(err, "INVALID_MAIL"), ShouldBeFalse)
	})
}
