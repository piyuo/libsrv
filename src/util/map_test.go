package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToString(t *testing.T) {
	assert := assert.New(t)
	m := map[string]string{
		"a": "1",
		"b": "2",
	}
	s := MapToString(m)
	assert.NotEmpty(s)

	m2 := MapFromString(s)
	assert.Equal("1", m2["a"])
	assert.Equal("2", m2["b"])
}

func TestMapToAndFromString(t *testing.T) {
	assert := assert.New(t)
	m := map[string]string{}
	s := MapToString(m)
	assert.Empty(s)

	m2 := MapFromString(s)
	assert.Equal(0, len(m2))
}

func TestMapFromString(t *testing.T) {
	assert := assert.New(t)
	m := MapFromString("=")
	assert.Equal(0, len(m))
}

func TestMapGetString(t *testing.T) {
	assert := assert.New(t)
	m := map[string]interface{}{"a": "b", "c": 1}
	assert.Equal("b", GetString(m, "a", "default"))
	assert.Equal("default", GetString(m, "b", "default"))
	assert.Equal("1", GetString(m, "c", "default"))
}

func TestMapGetFloat64(t *testing.T) {
	assert := assert.New(t)
	m := map[string]interface{}{"a": 123, "c": "b"}
	assert.Equal(float64(123), GetFloat64(m, "a", 1))
	assert.Equal(float64(0), GetFloat64(m, "b", 0))
	assert.Equal(float64(0), GetFloat64(m, "c", 0))
}
