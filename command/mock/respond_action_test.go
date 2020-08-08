package mock

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRespondAction(t *testing.T) {
	Convey("should execute the action and get response", t, func() {
		action := &RespondAction{}
		// action.Name = "hello"

		response, err := action.Do(context.Background())
		So(err, ShouldBeNil)
		So(response, ShouldNotBeNil)
		//sr := response.(*StringResponse)
		//So(sr.Text, ShouldEqual, "hi")
	})

}
