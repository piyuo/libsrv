package launch

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCheck(t *testing.T) {
	Convey("should pass check with no panic", t, func() {
		So(Checklist, ShouldNotPanic)

		os.Setenv("NAME", "")
		So(Checklist, ShouldPanic)
		os.Setenv("NAME", "not empty")

		os.Setenv("REGION", "")
		So(Checklist, ShouldPanic)
		os.Setenv("REGION", "not empty")

		os.Setenv("BRANCH", "")
		So(Checklist, ShouldPanic)
		os.Setenv("BRANCH", "not empty")

	})
}
