package token

import (
	"time"

	crypto "github.com/piyuo/libsrv/src/crypto"
	util "github.com/piyuo/libsrv/src/util"
	"github.com/pkg/errors"
)

// BaseToken implement Token
//
type BaseToken struct {
	Token

	// content is a key/value map
	//
	content map[string]string
}

//expiredFormat is expired time string format
//
const expiredFormat = "200601021504"

// keyExpired is expired key name
//
const keyExpired = "_"

// NewToken return a empty token
//
//	token := NewToken()
//
func NewToken() Token {
	return &BaseToken{
		content: map[string]string{},
	}
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

// FromString return Token from string or expired
//
//	token, expired, err := FromString(str)
//
func FromString(str string) (Token, bool, error) {
	everything, err := crypto.Decrypt(str)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to decrypt str:"+str)
	}
	content := util.MapFromString(everything)
	expired := content[keyExpired]
	if isExpired(expired) {
		return nil, true, nil
	}

	delete(content, keyExpired)
	return &BaseToken{
		content: content,
	}, false, nil
}

// ToString return string with expired time, after expired time the token will not read from string
//
//	expired := time.Now().UTC().Add(60 * time.Second)
//	str := token.ToString(expired)
//
func (c *BaseToken) ToString(expired time.Time) (string, error) {
	c.content[keyExpired] = expired.Format(expiredFormat)
	everything := util.MapToString(c.content)
	crypted, err := crypto.Encrypt(everything)
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt tokens")
	}
	return crypted, nil
}

// Get return value from key
//
//	value := token.Get("UserID")
//
func (c *BaseToken) Get(key string) string {
	return c.content[key]
}

// Set return value to key
//
//	token.Set("UserID","aa")
//
func (c *BaseToken) Set(key, value string) {
	c.content[key] = value
}

// Delete key
//
//	token.Delete("UserID")
//
func (c *BaseToken) Delete(key string) {
	delete(c.content, key)
}
