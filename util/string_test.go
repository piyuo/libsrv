package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindInString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("2", StringBetween("123", "1", "3"))
	assert.Equal("", StringBetween("123", "a", "3"))
	assert.Equal("", StringBetween("123", "1", "a"))
	assert.Equal("", StringBetween("111", "1", "1"))

	assert.Equal("1", StringBefore("123", "2"))
	assert.Equal("", StringBefore("123", "a"))
	assert.Equal("3", StringAfter("123", "2"))
	assert.Equal("", StringAfter("123", "a"))
	assert.Equal("", StringAfter("111", "1"))
}

func TestStringSplit(t *testing.T) {
	assert := assert.New(t)
	array := []string{"1", "2", "3"}
	str := ArrayToString(array)
	assert.NotEmpty(str)
	ary := StringToArray(str)
	assert.Equal(3, len(ary))
	assert.Equal("1", ary[0])
	assert.Equal("2", ary[1])
	assert.Equal("3", ary[2])
}

func TestGetHashcode(t *testing.T) {
	assert := assert.New(t)
	str := "hi"
	code := StringHash(str)
	assert.Greater(code, uint32(0))
	code2 := StringHash(str)
	assert.Equal(code, code2)
}
