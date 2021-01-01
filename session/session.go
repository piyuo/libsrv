package session

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/piyuo/libsrv/util"
)

// KeyContext define key used in ctx
//
type KeyContext int

const (
	// KeyRequest is context key name for request
	//
	KeyRequest KeyContext = iota

	// KeyUserID is context key name for user id
	//
	KeyUserID
)

// deadline cache os env COMMAND_DEADLINE value
//
var deadline time.Duration = -1

// SetDeadline set context deadline using os.Getenv("DEADLINE")
//
//	ctx = SetRequest(ctx,request)
//
func SetDeadline(ctx context.Context) (context.Context, context.CancelFunc) {

	if deadline == -1 {
		text := os.Getenv("DEADLINE")
		var err error
		ms, err := strconv.Atoi(text)
		if err != nil {
			ms = 20000
			fmt.Print("use default deadline 20 seconds")
		}
		deadline = time.Duration(ms)
	}
	expired := time.Now().Add(deadline * time.Millisecond)
	return context.WithDeadline(ctx, expired)
}

// GetRequest get current request from context
//
//	request := GetRequest(ctx)
//
func GetRequest(ctx context.Context) *http.Request {
	iRequest := ctx.Value(KeyRequest)
	if iRequest != nil {
		return iRequest.(*http.Request)
	}
	return nil
}

// SetRequest set request into ctx, this may used in log and data package
//
//	ctx = SetRequest(ctx,request)
//
func SetRequest(ctx context.Context, request *http.Request) context.Context {
	return context.WithValue(ctx, KeyRequest, request)
}

// SetUserID set UserID into ctx, this may used in log and data package
//
//	ctx = SetUserID(ctx,"user id")
//
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, KeyUserID, userID)
}

// GetUserID get current user id from context
//
//	userID := GetUserID(ctx)
//
func GetUserID(ctx context.Context) string {
	iUserID := ctx.Value(KeyUserID)
	if iUserID != nil {
		return iUserID.(string)
	}
	return ""
}

// GetIP return ip from current request, return empty if anything wrong
//
//	ip := GetIP(ctx)
//
func GetIP(ctx context.Context) string {
	value := ctx.Value(KeyRequest)
	if value == nil {
		return ""
	}
	req := value.(*http.Request)
	return util.GetIP(req)
}

// GetUserAgentID return short id from user agent. no version in here cause we used this for refresh token
//
//	ua := GetUserAgentID(ctx) // "iPhone,iOS,Safari"
//
func GetUserAgentID(ctx context.Context) string {
	value := ctx.Value(KeyRequest)
	if value == nil {
		return ""
	}
	req := value.(*http.Request)
	return util.GetUserAgentID(req)
}

// GetUserAgentString return short string with version info from user agent
//
//	ua := GetUserAgentString(ctx) // "iPhone,iOS 7.0,Safari 6.0"
//
func GetUserAgentString(ctx context.Context) string {
	value := ctx.Value(KeyRequest)
	if value == nil {
		return ""
	}
	req := value.(*http.Request)
	return util.GetUserAgentString(req)
}

// GetUserAgent return user agent from current request, return empty if anything wrong
//
//	ua := GetUserAgent(ctx) //"Mozilla/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X) AppleWebKit/546.10 (KHTML, like Gecko) Version/6.0 Mobile/7E18WD Safari/8536.25"
//
func GetUserAgent(ctx context.Context) string {
	value := ctx.Value(KeyRequest)
	if value == nil {
		return ""
	}
	req := value.(*http.Request)
	return util.GetUserAgent(req)
}

// GetLocale return locale from current request, return en-us if anything else
//
//	lang := GetLocale(ctx)
//
func GetLocale(ctx context.Context) string {
	value := ctx.Value(KeyRequest)
	if value == nil {
		return "en_US"
	}
	req := value.(*http.Request)
	return util.GetLocale(req)
}
