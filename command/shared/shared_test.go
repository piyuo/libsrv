package shared

import (
	"context"
	"testing"

	sharedcommands "github.com/piyuo/go-libsrv/command/shared/commands"

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
		text := Text("hi").(*sharedcommands.Text)
		So(text.Value, ShouldEqual, "hi")
	})

	Convey("should create number response'", t, func() {
		num := Number(201).(*sharedcommands.Num)
		So(num.Value, ShouldEqual, 201)
	})

	Convey("should create error response'", t, func() {
		err := Error(ErrorUnknown, "tag").(*sharedcommands.Err)
		So(err.Code, ShouldEqual, int32(ErrorUnknown))
		So(err.Tag, ShouldEqual, "tag")
	})

}
