package mock

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDeadlineAction(t *testing.T) {
	Convey("should execute the action and get response", t, func() {
		action := &DeadlineAction{}
		// action.Name = "hello"

		response, err := action.Do(context.Background())
		So(err, ShouldNotBeNil)
		So(response, ShouldBeNil)
	})

}
