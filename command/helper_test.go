package command

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWriteResponse(t *testing.T) {
	Convey("should write binary", t, func() {
		w := httptest.NewRecorder()
		bytes := newTestAction(textLong)
		writeBinary(w, bytes)
		writeText(w, "code")
		writeError(w, errors.New("error"), 500, "error")
		writeBadRequest(context.Background(), w, "message")
	})
}


func TestIsSlow(t *testing.T) {
	Convey("should determine slow work", t, func() {
		// 3 seconds execution time is not slow
		So(IsSlow(5000), ShouldEqual, 0)
		// 20 seconds execution time is really slow
		So(IsSlow(20000000), ShouldBeGreaterThan, 5000)
	})
}
