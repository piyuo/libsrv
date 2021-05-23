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

func TestResourceKey(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", "zh-tw")
	ctx := context.WithValue(context.Background(), env.KeyContextRequest, req)
	assert.Equal("name_zh_TW.json", ResourceKey(ctx, "name", ".json"))
}

func TestResourcePath(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", "zh-tw")
	ctx := context.WithValue(context.Background(), env.KeyContextRequest, req)
	assert.Equal("assets/i18n/name_zh_TW.json", ResourcePath(ctx, "name", ".json"))
}

func TestResource(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", "en_US")
	ctx := context.WithValue(context.Background(), env.KeyContextRequest, req)
	json, err := Resource(ctx, "mock", ".json")
	assert.Nil(err)
	assert.Equal("world", json["hello"])
	//get from cache
	json, err = Resource(ctx, "mock", ".json")
	assert.Nil(err)
	assert.Equal("world", json["hello"])
}

func TestResourceWithoutCache(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", "en_US")
	ctx := context.WithValue(context.Background(), env.KeyContextRequest, req)
	json, bytes, err := ResourceWithoutCache(ctx, "mock", ".json")
	assert.Nil(err)
	assert.NotNil(bytes)
	assert.Equal("world", json["hello"])
}

func TestResourceNotFound(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", "en_US")
	ctx := context.WithValue(context.Background(), env.KeyContextRequest, req)
	json, err := Resource(ctx, "notExist", ".json")
	assert.NotNil(err)
	assert.Nil(json)
}
