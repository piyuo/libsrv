package i18n

import (
	"context"
	"net/http"
	"testing"

	"github.com/piyuo/libsrv/env"
	"github.com/stretchr/testify/assert"
)

func TestIsPredefined(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	exist, predefine := IsPredefined("en-us")
	assert.True(exist)
	assert.Equal("en_US", predefine)

	exist, predefine = IsPredefined("zh-tw")
	assert.True(exist)
	assert.Equal("zh_TW", predefine)

	exist, predefine = IsPredefined("en")
	assert.False(exist)
	assert.Empty(predefine)
}

func TestAcceptLanguage(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	locale := acceptLanguage("")
	assert.Equal("en_US", locale)

	locale = acceptLanguage("en-us")
	assert.Equal("en_US", locale)

	locale = acceptLanguage("zh_TW")
	assert.Equal("zh_TW", locale)

	locale = acceptLanguage("da, en-us;q=0.8, en;q=0.7")
	assert.Equal("en_US", locale)

	locale = acceptLanguage("da, zh-cn;q=0.8, zh-tw;q=0.7")
	assert.Equal("zh_CN", locale)
}

func TestGetLocaleFromRequest(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/whatever", nil)
	req.Header.Add("Accept-Language", "da, en-xx;q=0.8, en;q=0.7")
	locale := GetLocaleFromRequest(req)
	assert.Equal("en_US", locale) //nothing match, use en_US

	req.Header = map[string][]string{}
	//empty accept-language will result en-us
	req.Header.Add("Accept-Language", "")
	locale = GetLocaleFromRequest(req)
	assert.Equal("en_US", locale)

	req.Header = map[string][]string{}
	//will convert accept language to predefined locale
	req.Header.Add("Accept-Language", "EN-US")
	locale = GetLocaleFromRequest(req)
	assert.Equal("en_US", locale)
	req.Header = map[string][]string{}
}

func TestGetLocaleFromContext(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := context.Background()
	assert.Equal("en_US", GetLocaleFromContext(ctx))

	req.Header.Add("Accept-Language", "zh-cn")
	ctx = context.WithValue(context.Background(), env.KeyContextRequest, req)
	assert.Equal("zh_CN", GetLocaleFromContext(ctx))
}

func TestLocaleFilename(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", "zh-tw")
	ctx := context.WithValue(context.Background(), env.KeyContextRequest, req)
	assert.Equal("name_zh_TW.json", LocaleFilename(ctx, "name", ".json"))
}

func TestJSON(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", "en_US")
	ctx := context.WithValue(context.Background(), env.KeyContextRequest, req)
	json, err := JSON(ctx, "mock", ".json", 0)
	assert.Nil(err)
	assert.Equal("world", json["hello"])
	//get from cache
	json, err = JSON(ctx, "mock", ".json", 0)
	assert.Nil(err)
	assert.Equal("world", json["hello"])
}

func TestResourceNotFound(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", "en_US")
	ctx := context.WithValue(context.Background(), env.KeyContextRequest, req)
	json, err := JSON(ctx, "notExist", ".json", 0)
	assert.NotNil(err)
	assert.Nil(json)
}

func TestText(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", "en_US")
	ctx := context.WithValue(context.Background(), env.KeyContextRequest, req)
	txt, err := Text(ctx, "mock", ".json", 0)
	assert.Nil(err)
	assert.NotEmpty(txt)
	//get from cache
	txt, err = Text(ctx, "mock", ".json", 0)
	assert.Nil(err)
	assert.NotEmpty(txt)
}
