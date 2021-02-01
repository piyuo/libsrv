package env

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIPLocale(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(GetRequest(ctx))

	req, _ := http.NewRequest("GET", "/", nil)
	ctx = SetRequest(ctx, req)
	assert.NotNil(ctx)

	req2 := GetRequest(ctx)
	assert.Equal(req, req2)
}

func TestGetIPAndLocale(t *testing.T) {
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := context.Background()
	assert.Empty(GetIP(ctx))

	req.Header.Add("Accept-Language", "zh-cn")
	req.RemoteAddr = "[::1]:80"
	ctx = context.WithValue(context.Background(), KeyContextRequest, req)
	assert.Equal("::1", GetIP(ctx))
}

func TestUserAgent(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	str := GetUserAgentString(ctx)
	assert.Empty(str)

	req, _ := http.NewRequest("GET", "/whatever", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X) AppleWebKit/546.10 (KHTML, like Gecko) Version/6.0 Mobile/7E18WD Safari/8536.25")

	assert.Empty(GetUserAgent(ctx))
	assert.Empty(GetUserAgentID(ctx))

	ctx = context.WithValue(context.Background(), KeyContextRequest, req)

	ua := GetUserAgent(ctx)
	assert.NotEmpty(ua)
	id := GetUserAgentID(ctx)
	assert.Equal("iPhone,iOS,Safari", id)
	str = GetUserAgentString(ctx)
	assert.Equal("iPhone,iOS 7.0,Safari 6.0", str)
}

func TestUserID(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	userID := GetUserID(ctx)
	assert.Empty(userID)
	ctx = SetUserID(ctx, "id")
	userID = GetUserID(ctx)
	assert.Equal("id", userID)
}

func TestAccountID(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	accountID := GetAccountID(ctx)
	assert.Empty(accountID)
	ctx = SetAccountID(ctx, "id")
	accountID = GetAccountID(ctx)
	assert.Equal("id", accountID)
}
