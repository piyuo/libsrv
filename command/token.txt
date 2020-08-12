package command

import (
	"context"
	"net/http"
	"time"

	crypto "github.com/piyuo/libsrv/crypto"
	util "github.com/piyuo/libsrv/util"
	"github.com/pkg/errors"
)

//ErrTokenRequired mean service need access  token
//
var ErrTokenRequired = errors.New("TOKEN_REQUIRED")

// cookieKey is token key name in cookie
//
const cookieKey = "T"

// expired is token expired time
//
const expired = 3600 // expired in one hour

//expiredFormat is expired time string format
//
const expiredFormat = "200601021504"

// getExpired return expired time in string format,expired in seconds
//
//	txt := getExpired(300) // 5 min
//	So(txt, ShouldNotBeEmpty)
//
func getExpired(expired int) string {
	t := time.Now().UTC().Add(time.Second * time.Duration(expired))
	return t.Format(expiredFormat)
}

//	isExpired check string format datetime is expired, return true if anything wrong
//
//	expired = isExpired("200001010101")
//	So(expired, ShouldBeTrue)
//
func isExpired(str string) bool {
	expired, err := time.Parse(expiredFormat, str)
	if err != nil {
		return true
	}
	if expired.After(time.Now()) {
		return false
	}
	return true
}

// TokenRequired check is token is exist and has id, return error if token not exist
//
//	err := TokenRequired(ctx)
//
func TokenRequired(ctx context.Context) error {
	value := ctx.Value(keyToken)
	if value == nil {
		return ErrTokenRequired
	}

	m := value.(map[string]string)
	if m["id"] == "" {
		return ErrTokenRequired
	}
	return nil
}

// Tokens return token map in context, if no map in context return new map instead
//
//	m := Tokens(ctx)
//
func Tokens(ctx context.Context) map[string]string {
	// no need to check map exist in context, server will always place token map in context
	value := ctx.Value(keyToken)
	if value == nil {
		return map[string]string{}
	}
	return value.(map[string]string)
}

// SetToken set token to persist through cookie until expired, token["id"] will be used in log
//
//	err := SetToken("id","111-222)
//
func SetToken(ctx context.Context, key, value string) {
	tokens := Tokens(ctx)
	tokens[key] = value
}

// GetToken get token from context
//
//	value := GetToken(ctx,"id")
//
func GetToken(ctx context.Context, key string) string {
	tokens := Tokens(ctx)
	return tokens[key]
}

// contextToCookie persist context token through cookie
//
//	err := contextToCookie(w)
//
func contextToCookie(ctx context.Context, w http.ResponseWriter) error {
	tokens := Tokens(ctx)
	if len(tokens) == 0 {
		//remove cookie
		http.SetCookie(w, &http.Cookie{
			Name:     cookieKey,
			Value:    "",
			HttpOnly: true,
			MaxAge:   0,
		})
		return nil
	}

	//reset expired time, even nothing in token is change. this make sure expired time reset every command call
	tokens["expired"] = getExpired(expired)
	everything := util.MapToString(tokens)
	crypted, err := crypto.Encrypt(everything)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt tokens")
	}
	http.SetCookie(w, &http.Cookie{
		Name:     cookieKey,
		Value:    crypted,
		HttpOnly: true,
		Secure:   true,
		MaxAge:   expired,
	})
	return nil
}

// cookieToToken reload token from cookie
//
//	err := cookieToToken(ctx,w)
//
func contextFromCookie(ctx context.Context, r *http.Request) (context.Context, error) {
	if len(r.Cookies()) == 0 {
		return context.WithValue(ctx, keyToken, map[string]string{}), nil
	}

	cookie, err := r.Cookie(cookieKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get context token from cookie[\""+cookieKey+"\"]")
	}
	if cookie.Value == "" {
		return context.WithValue(ctx, keyToken, map[string]string{}), nil
	}
	everything, err := crypto.Decrypt(cookie.Value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt text("+cookie.Value+")")
	}
	tokens := util.MapFromString(everything)
	expired := tokens["expired"]
	if isExpired(expired) {
		return context.WithValue(ctx, keyToken, map[string]string{}), nil
	}
	delete(tokens, "expired")
	return context.WithValue(ctx, keyToken, tokens), nil
}
