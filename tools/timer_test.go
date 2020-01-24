package tools

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTimer(t *testing.T) {
	Convey("should use timer'", t, func() {
		timer := NewTimer()
		timer.Start()
		time.Sleep(1 * time.Millisecond)
		ms := timer.Stop()
		So(ms >= 1, ShouldBeTrue)
	})
}
