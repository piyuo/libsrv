package util

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUserAgent(t *testing.T) {
	Convey("should get useragent", t, func() {
		req, _ := http.NewRequest("GET", "/whatever", nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X) AppleWebKit/546.10 (KHTML, like Gecko) Version/6.0 Mobile/7E18WD Safari/8536.25")

		ua := GetUserAgent(req)
		So(ua, ShouldNotBeEmpty)

		browserName, browserVer, osName, osVer, device := ParseUserAgent(ua)
		So(browserName, ShouldNotBeEmpty)
		So(browserVer, ShouldNotBeEmpty)
		So(osName, ShouldNotBeEmpty)
		So(osVer, ShouldNotBeEmpty)
		So(device, ShouldNotBeEmpty)

		id := GetUserAgentID(req)
		So(id, ShouldEqual, "iPhone,iOS,Safari")

		str := GetUserAgentString(req)
		So(str, ShouldEqual, "iPhone,iOS 7.0,Safari 6.0")
	})
}

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
