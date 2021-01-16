package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldConvertMapToString(t *testing.T) {
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

func TestShouldAllowEmptyMap(t *testing.T) {
	assert := assert.New(t)
	m := map[string]string{}
	s := MapToString(m)
	assert.Empty(s)

	m2 := MapFromString(s)
	assert.Equal(0, len(m2))
}

func TestShouldReturnEmptyMap(t *testing.T) {
	assert := assert.New(t)
	m := MapFromString("=")
	assert.Equal(0, len(m))
}
