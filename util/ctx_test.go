package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCanceledCtx(t *testing.T) {
	Convey("should return canceled ctx", t, func() {
		ctx := CanceledCtx()
		So(ctx.Err(), ShouldNotBeNil)
	})
}
