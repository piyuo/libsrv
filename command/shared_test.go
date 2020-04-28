package command

import (
	"context"
	"testing"

	shared "github.com/piyuo/libsrv/command/shared"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShared(t *testing.T) {
	Convey("should get error from empty context'", t, func() {
		ctx := context.Background()
		token, response := Token(ctx)
		So(token, ShouldBeNil)
		So(response, ShouldNotBeNil)
	})

	Convey("should create text response'", t, func() {
		text := String("hi").(*shared.Text)
		So(text.Value, ShouldEqual, "hi")
	})

	Convey("should create number response'", t, func() {
		num := Number(201).(*shared.Num)
		So(num.Value, ShouldEqual, 201)
	})

	Convey("should create error response'", t, func() {
		err := Error("errCode").(*shared.Err)
		So(err.Code, ShouldEqual, "errCode")
	})

}
