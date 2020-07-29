package command

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExpired(t *testing.T) {
	Convey("should check expired date", t, func() {
		txt := getExpired(300) // 5 min
		So(txt, ShouldNotBeEmpty)

		expired := isExpired(txt)
		So(expired, ShouldBeFalse)

		expired = isExpired("200001010101")
		So(expired, ShouldBeTrue)

		expired = isExpired("300001010101")
		So(expired, ShouldBeFalse)
	})
}

func TestTokens(t *testing.T) {
	Convey("should get tokens from context", t, func() {
		ctx := context.Background()
		m := Tokens(ctx)
		So(m, ShouldNotBeNil)
		So(m["not-exist"], ShouldBeEmpty)

		ctx = context.WithValue(ctx, keyToken, map[string]string{"a": "1"})
		m = Tokens(ctx)
		So(m, ShouldNotBeNil)
		So(m["a"], ShouldEqual, "1")

		m["a"] = "2"
		m2 := Tokens(ctx)
		So(m2, ShouldNotBeNil)
		So(m2["a"], ShouldEqual, "2")
	})
}

func TestToken(t *testing.T) {
	Convey("should get/set token", t, func() {
		ctx := context.Background()
		ctx = context.WithValue(ctx, keyToken, map[string]string{"a": "1"})
		value := GetToken(ctx, "a")
		So(value, ShouldEqual, "1")
		value = GetToken(ctx, "b")
		So(value, ShouldEqual, "")

		SetToken(ctx, "a", "2")
		value = GetToken(ctx, "a")
		So(value, ShouldEqual, "2")
	})
}

func TestEmptyContextToken(t *testing.T) {
	Convey("should handle empty context", t, func() {
		req, _ := http.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		ctx := context.Background()
		ctx = context.WithValue(ctx, keyToken, map[string]string{})

		err := contextToCookie(ctx, rec)
		So(err, ShouldBeNil)

		resp := rec.Result()
		So(resp, ShouldNotBeNil)
		cookie := resp.Cookies()[0]
		So(cookie.Name, ShouldEqual, cookieKey)
		So(cookie.Value, ShouldBeEmpty)

		ctx, err = contextFromCookie(context.Background(), req)
		So(err, ShouldBeNil)
		tokens := Tokens(ctx)
		So(len(tokens), ShouldEqual, 0)
	})
}

func TestContextToken(t *testing.T) {
	Convey("should to/from cookie", t, func() {
		rec := httptest.NewRecorder()
		ctx := context.Background()
		ctx = context.WithValue(ctx, keyToken, map[string]string{})

		SetToken(ctx, "a", "1")
		err := contextToCookie(ctx, rec)
		So(err, ShouldBeNil)

		resp := rec.Result()
		So(resp, ShouldNotBeNil)
		cookie := resp.Cookies()[0]
		So(cookie.Name, ShouldEqual, cookieKey)
		So(cookie.Value, ShouldNotBeEmpty)

		req, _ := http.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: cookieKey, Value: cookie.Value})
		ctx, err = contextFromCookie(context.Background(), req)
		So(err, ShouldBeNil)
		value := GetToken(ctx, "a")
		So(value, ShouldEqual, "1")
	})
}

func TestErrorCookie(t *testing.T) {
	Convey("should to/from cookie", t, func() {

		req, _ := http.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: cookieKey, Value: "abc"}) //someone try to modify cookie
		ctx, err := contextFromCookie(context.Background(), req)
		So(err, ShouldNotBeNil)

		req, _ = http.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: cookieKey, Value: ""}) //someone try to modify cookie
		ctx, err = contextFromCookie(context.Background(), req)
		So(err, ShouldBeNil)
		tokens := Tokens(ctx)
		So(len(tokens), ShouldEqual, 0)
	})
}
