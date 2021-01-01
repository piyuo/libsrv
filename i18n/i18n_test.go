package i18n

import (
	"context"
	"net/http"
	"testing"

	"github.com/piyuo/libsrv/session"
	. "github.com/smartystreets/goconvey/convey"
)

func TestIsPredefined(t *testing.T) {
	Convey("should return predefined locale", t, func() {

		exist, predefine := IsPredefined("en-us")
		So(exist, ShouldBeTrue)
		So(predefine, ShouldEqual, "en_US")

		exist, predefine = IsPredefined("zh-tw")
		So(exist, ShouldBeTrue)
		So(predefine, ShouldEqual, "zh_TW")

		exist, predefine = IsPredefined("en")
		So(exist, ShouldBeFalse)
		So(predefine, ShouldBeEmpty)

	})
}

func TestAcceptLanguage(t *testing.T) {
	Convey("should get accept language", t, func() {
		locale := acceptLanguage("")
		So(locale, ShouldEqual, "en_US")

		locale = acceptLanguage("en-us")
		So(locale, ShouldEqual, "en_US")

		locale = acceptLanguage("zh_TW")
		So(locale, ShouldEqual, "zh_TW")

		locale = acceptLanguage("da, en-us;q=0.8, en;q=0.7")
		So(locale, ShouldEqual, "en_US")

		locale = acceptLanguage("da, zh-cn;q=0.8, zh-tw;q=0.7")
		So(locale, ShouldEqual, "zh_CN")
	})
}

func TestGetLocaleFromRequest(t *testing.T) {
	Convey("should get accept language", t, func() {
		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Add("Accept-Language", "da, en-xx;q=0.8, en;q=0.7")
		locale := GetLocaleFromRequest(req)
		So(locale, ShouldEqual, "en_US") //nothing match, use en_US

		req.Header = map[string][]string{}
		//empty accept-language will result en-us
		req.Header.Add("Accept-Language", "")
		locale = GetLocaleFromRequest(req)
		So(locale, ShouldEqual, "en_US")

		req.Header = map[string][]string{}
		//will convert accept language to predefined locale
		req.Header.Add("Accept-Language", "EN-US")
		locale = GetLocaleFromRequest(req)
		So(locale, ShouldEqual, "en_US")
		req.Header = map[string][]string{}

	})
}

func TestGetLocaleFromContext(t *testing.T) {
	Convey("should get ip and locale", t, func() {
		req, _ := http.NewRequest("GET", "/", nil)
		ctx := context.Background()
		So(GetLocaleFromContext(ctx), ShouldEqual, "en_US")

		req.Header.Add("Accept-Language", "zh-cn")
		ctx = context.WithValue(context.Background(), session.KeyRequest, req)
		So(GetLocaleFromContext(ctx), ShouldEqual, "zh_CN")
	})
}

func TestResourceKey(t *testing.T) {
	Convey("should get resource key", t, func() {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Add("Accept-Language", "zh-tw")
		ctx := context.WithValue(context.Background(), session.KeyRequest, req)
		So(ResourceKey(ctx, "name"), ShouldEqual, "name_zh_TW")
	})
}

func TestResourcePath(t *testing.T) {
	Convey("should get resource path", t, func() {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Add("Accept-Language", "zh-tw")
		ctx := context.WithValue(context.Background(), session.KeyRequest, req)
		So(ResourcePath(ctx, "name"), ShouldEqual, "assets/i18n/name_zh_TW.json")
	})
}

func TestResource(t *testing.T) {
	Convey("should get resource path", t, func() {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Add("Accept-Language", "en_US")
		ctx := context.WithValue(context.Background(), session.KeyRequest, req)
		json, err := Resource(ctx, "mock")
		So(err, ShouldBeNil)
		So(json["hello"], ShouldEqual, "world")
	})
}

func TestResourceNotFound(t *testing.T) {
	Convey("should get resource path", t, func() {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Add("Accept-Language", "en_US")
		ctx := context.WithValue(context.Background(), session.KeyRequest, req)
		json, err := Resource(ctx, "notExist")
		So(err, ShouldNotBeNil)
		So(json, ShouldBeNil)
	})
}
