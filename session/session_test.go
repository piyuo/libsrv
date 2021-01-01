package session

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDeadline(t *testing.T) {
	Convey("should return deadline", t, func() {
		ctx := context.Background()
		So(ctx.Err(), ShouldBeNil)

		backup := os.Getenv("DEADLINE")
		os.Setenv("DEADLINE", "20")
		ctx, cancel := SetDeadline(ctx)
		defer cancel()

		So(ctx.Err(), ShouldBeNil)
		time.Sleep(time.Duration(31) * time.Millisecond)
		So(ctx.Err(), ShouldNotBeNil)

		deadline = -1 // remove cache
		os.Setenv("DEADLINE", backup)
	})
}

func TestDeadlineNotSet(t *testing.T) {
	Convey("should return deadline", t, func() {
		ctx := context.Background()
		So(ctx.Err(), ShouldBeNil)

		backup := os.Getenv("DEADLINE")
		os.Setenv("DEADLINE", "")
		ctx, cancel := SetDeadline(ctx)
		defer cancel()

		time.Sleep(time.Duration(21) * time.Millisecond)
		So(ctx.Err(), ShouldBeNil) // default expired is in 20,000ms
		deadline = -1              // remove cache
		os.Setenv("DEADLINE", backup)
	})
}

func TestRequest(t *testing.T) {
	Convey("should get ip and locale", t, func() {
		ctx := context.Background()
		So(GetRequest(ctx), ShouldBeNil)

		req, _ := http.NewRequest("GET", "/", nil)
		ctx = SetRequest(ctx, req)
		So(ctx, ShouldNotBeNil)

		req2 := GetRequest(ctx)
		So(req, ShouldEqual, req2)
	})
}

func TestGetIPAndLocale(t *testing.T) {
	Convey("should get ip and locale", t, func() {
		req, _ := http.NewRequest("GET", "/", nil)
		ctx := context.Background()
		So(GetIP(ctx), ShouldEqual, "")
		So(GetLocale(ctx), ShouldEqual, "en_US")

		req.Header.Add("Accept-Language", "zh-cn")
		req.RemoteAddr = "[::1]:80"
		ctx = context.WithValue(context.Background(), KeyRequest, req)
		So(GetIP(ctx), ShouldEqual, "::1")
		So(GetLocale(ctx), ShouldEqual, "zh_CN")
	})
}

func TestUserAgent(t *testing.T) {
	Convey("should get useragent", t, func() {
		ctx := context.Background()
		str := GetUserAgentString(ctx)
		So(str, ShouldBeEmpty)

		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X) AppleWebKit/546.10 (KHTML, like Gecko) Version/6.0 Mobile/7E18WD Safari/8536.25")

		So(GetUserAgent(ctx), ShouldEqual, "")
		So(GetUserAgentID(ctx), ShouldEqual, "")

		ctx = context.WithValue(context.Background(), KeyRequest, req)

		ua := GetUserAgent(ctx)
		So(ua, ShouldNotBeEmpty)
		id := GetUserAgentID(ctx)
		So(id, ShouldEqual, "iPhone,iOS,Safari")
		str = GetUserAgentString(ctx)
		So(str, ShouldEqual, "iPhone,iOS 7.0,Safari 6.0")
	})
}

func TestUserID(t *testing.T) {
	Convey("should get/set user id in context", t, func() {
		ctx := context.Background()
		userID := GetUserID(ctx)
		So(userID, ShouldBeEmpty)
		ctx = SetUserID(ctx, "id")
		userID = GetUserID(ctx)
		So(userID, ShouldEqual, "id")
	})
}
