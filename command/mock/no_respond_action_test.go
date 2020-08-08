package mock

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNoRespondAction(t *testing.T) {
	Convey("should execute the action and get response", t, func() {
		action := &NoRespondAction{}
		// action.Name = "hello"

		response, err := action.Do(context.Background())
		So(err, ShouldBeNil)
		So(response, ShouldBeNil)
		//sr := response.(*StringResponse)
		//So(sr.Text, ShouldEqual, "hi")
	})

}
