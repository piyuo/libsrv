package command

import (
	"testing"

	shared "github.com/piyuo/libsrv/command/shared"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShared(t *testing.T) {

	Convey("should create text response", t, func() {
		text := String("hi").(*shared.Text)
		So(text.Value, ShouldEqual, "hi")
	})

	Convey("should create number response", t, func() {
		num := Number(201).(*shared.Num)
		So(num.Value, ShouldEqual, 201)
	})

	Convey("should create error response", t, func() {
		err := Error("errCode").(*shared.Err)
		So(err.Code, ShouldEqual, "errCode")
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
