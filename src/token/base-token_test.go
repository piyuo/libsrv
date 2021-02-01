package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpiredToken(t *testing.T) {
	assert := assert.New(t)
	tt := time.Now().UTC().Add(60 * time.Second)
	txt := tt.Format(expiredFormat)

	expired := isExpired(txt)
	assert.False(expired)

	expired = isExpired("200001010101")
	assert.True(expired)

	expired = isExpired("300001010101")
	assert.False(expired)
}

func TestGetSetDeleteToken(t *testing.T) {
	assert := assert.New(t)

	token := NewToken()

	value := token.Get("a")
	assert.Equal("", value)

	token.Set("a", "1")
	value = token.Get("a")
	assert.Equal("1", value)

	token.Delete("a")
	value = token.Get("a")
	assert.Equal("", value)
}

func TestTokenFromToString(t *testing.T) {
	assert := assert.New(t)
	token := NewToken()
	token.Set("a", "1")
	expired := time.Now().UTC().Add(60 * time.Second)
	crypted, err := token.ToString(expired)
	assert.Nil(err)
	assert.NotEmpty(crypted)

	token2, isExpired, err := FromString(crypted)
	assert.Nil(err)
	assert.False(isExpired)

	value := token2.Get("a")
	assert.Equal("1", value)
}

func TestTokenExpired(t *testing.T) {
	assert := assert.New(t)
	token := NewToken()
	token.Set("a", "1")

	expired := time.Now().UTC().Add(-60 * time.Second)
	crypted, err := token.ToString(expired)
	assert.Nil(err)
	assert.NotEmpty(crypted)

	token2, isExpired, err := FromString(crypted)
	assert.Nil(err)
	assert.True(isExpired)
	assert.Nil(token2)
}

func TestInvalidToken(t *testing.T) {
	assert := assert.New(t)
	token, expired, err := FromString("")
	assert.NotNil(err)
	assert.False(expired)
	assert.Nil(token)

	token, expired, err = FromString("123213123")
	assert.NotNil(err)
	assert.False(expired)
	assert.Nil(token)
}

func TestIsExpired(t *testing.T) {
	assert := assert.New(t)
	result := isExpired("a")
	assert.True(result)
}
