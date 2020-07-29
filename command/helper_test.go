package command

import (
	"context"
	"net/http"
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

func TestDeadline(t *testing.T) {
	Convey("should return deadline", t, func() {
		dateline := getDeadline()
		So(dateline.After(time.Now()), ShouldBeTrue)

		//dateline should not greater than 10 min.
		tenMinutesLater := time.Now().Add(24 * time.Hour)
		So(dateline.Before(tenMinutesLater), ShouldBeTrue)
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

func TestGetIPAndLocale(t *testing.T) {
	Convey("should get ip and locale", t, func() {
		req, _ := http.NewRequest("GET", "/", nil)
		ctx := context.Background()
		So(GetIP(ctx), ShouldEqual, "")
		So(GetLocale(ctx), ShouldEqual, "en-us")

		req.Header.Add("Accept-Language", "zh-cn")
		req.RemoteAddr = "[::1]:80"
		ctx = context.WithValue(context.Background(), keyRequest, req)
		So(GetIP(ctx), ShouldEqual, "::1")
		So(GetLocale(ctx), ShouldEqual, "zh-cn")
	})
}
