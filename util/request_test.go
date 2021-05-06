package util

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldGetUserAgent(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/whatever", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X) AppleWebKit/546.10 (KHTML, like Gecko) Version/6.0 Mobile/7E18WD Safari/8536.25")

	ua := GetUserAgent(req)
	assert.NotEmpty(ua)

	browserName, browserVer, osName, osVer, device := ParseUserAgent(ua)
	assert.NotEmpty(browserName)
	assert.NotEmpty(browserVer)
	assert.NotEmpty(osName)
	assert.NotEmpty(osVer)
	assert.NotEmpty(device)

	id := GetUserAgentID(req)
	assert.Equal("iPhone,iOS,Safari", id)

	str := GetUserAgentString(req)
	assert.Equal("iPhone,iOS 7.0,Safari 6.0", str)
}

func TestShouldGetIP(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/whatever", nil)
	ip := GetIP(req)
	assert.Empty(ip)

	req.RemoteAddr = "[::1]:80"
	assert.Equal("::1", GetIP(req))
	req.RemoteAddr = ""

	//wrong remote addr
	req.RemoteAddr = "xxx"
	assert.Equal("", GetIP(req))
	req.RemoteAddr = ""

	req.Header.Add("X-Real-IP", "12.34.56.78")
	assert.Equal("12.34.56.78", GetIP(req))
	req.Header = map[string][]string{}

	req.Header.Add("X-Forwarded-For", "23.45.67.89,12.34.56.78")
	assert.Equal("23.45.67.89", GetIP(req))
	req.Header = map[string][]string{}
}
