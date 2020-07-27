package util

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetIP(t *testing.T) {
	Convey("should get ip", t, func() {
		req, _ := http.NewRequest("GET", "/whatever", nil)
		ip := GetIP(req)
		So(ip, ShouldBeEmpty)

		req.RemoteAddr = "[::1]:80"
		So(GetIP(req), ShouldEqual, "::1")
		req.RemoteAddr = ""

		//wrong remote addr
		req.RemoteAddr = "xxx"
		So(GetIP(req), ShouldEqual, "")
		req.RemoteAddr = ""

		req.Header.Add("X-Real-IP", "12.34.56.78")
		So(GetIP(req), ShouldEqual, "12.34.56.78")
		req.Header = map[string][]string{}

		req.Header.Add("X-Forwarded-For", "23.45.67.89,12.34.56.78")
		So(GetIP(req), ShouldEqual, "23.45.67.89")
		req.Header = map[string][]string{}

	})
}

func TestAcceptLanguage(t *testing.T) {
	Convey("should get accept language", t, func() {
		list := parseAcceptLanguage("")
		So(len(list), ShouldEqual, 1)
		So(list[0], ShouldEqual, "en-us")

		list = parseAcceptLanguage("da, en-gb;q=0.8, en;q=0.7")
		So(list[0], ShouldEqual, "da")
		So(list[1], ShouldEqual, "en-gb")
		So(list[2], ShouldEqual, "en")

		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Add("Accept-Language", "da, en-gb;q=0.8, en;q=0.7")
		list = GetAcceptLanguage(req)
		So(list[0], ShouldEqual, "da")
		So(list[1], ShouldEqual, "en-gb")
		So(list[2], ShouldEqual, "en")
		req.Header = map[string][]string{}

		req.Header.Add("Accept-Language", "da, en-gb;q=0.8, en;q=0.7")
		lang := GetLanguage(req)
		So(lang, ShouldEqual, "da")
		req.Header = map[string][]string{}

		//empty accept-language will result en-us
		req.Header.Add("Accept-Language", "")
		lang = GetLanguage(req)
		So(lang, ShouldEqual, "en-us")
		req.Header = map[string][]string{}

		//language will be lower case
		req.Header.Add("Accept-Language", "EN-US")
		lang = GetLanguage(req)
		So(lang, ShouldEqual, "en-us")
		req.Header = map[string][]string{}

	})
}
