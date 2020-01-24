package tools

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUUID(t *testing.T) {
	Convey("should generate uuid", t, func() {
		id := UUID()
		So(id, ShouldNotBeEmpty)
	})

}
