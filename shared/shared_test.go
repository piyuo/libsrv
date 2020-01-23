package shared

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShared(t *testing.T) {
	Convey("should generate log head'", t, func() {
		ctx := context.Background()
		token, response := NeedToken(ctx)
		So(token, ShouldBeNil)
		So(response, ShouldNotBeNil)
	})
}
