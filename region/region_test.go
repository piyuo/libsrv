package region

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegion(t *testing.T) {
	Convey("should get current region", t, func() {
		So(Current, ShouldNotBeEmpty)
	})
}
