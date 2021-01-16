package session

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDeadline(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("DEADLINE")
	os.Setenv("DEADLINE", "20")
	ctx, cancel := SetDeadline(ctx)
	defer cancel()

	assert.Nil(ctx.Err())
	time.Sleep(time.Duration(31) * time.Millisecond)
	assert.NotNil(ctx.Err())

	deadline = -1 // remove cache
	os.Setenv("DEADLINE", backup)
}

func TestDeadlineNotSet(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	assert.Nil(ctx.Err())

	backup := os.Getenv("DEADLINE")
	os.Setenv("DEADLINE", "")
	ctx, cancel := SetDeadline(ctx)
	defer cancel()

	time.Sleep(time.Duration(21) * time.Millisecond)
	assert.Nil(ctx.Err()) // default expired is in 20,000ms
	deadline = -1         // remove cache
	os.Setenv("DEADLINE", backup)
}

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
	ctx = context.WithValue(context.Background(), KeyRequest, req)
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

	ctx = context.WithValue(context.Background(), KeyRequest, req)

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
