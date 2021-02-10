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
	assert.Equal("b", MapGetString(m, "a", "default"))
	assert.Equal("default", MapGetString(m, "b", "default"))
	assert.Equal("1", MapGetString(m, "c", "default"))
}

func TestMapGetFloat64(t *testing.T) {
	assert := assert.New(t)
	m := map[string]interface{}{"a": 123, "c": "b"}
	assert.Equal(float64(123), MapGetFloat64(m, "a", 1))
	assert.Equal(float64(0), MapGetFloat64(m, "b", 0))
	assert.Equal(float64(0), MapGetFloat64(m, "c", 0))

	f64 := map[string]interface{}{"a": float64(123)}
	assert.Equal(float64(123), MapGetFloat64(f64, "a", 1))
	f32 := map[string]interface{}{"a": float32(123)}
	assert.Equal(float64(123), MapGetFloat64(f32, "a", 1))
	i64 := map[string]interface{}{"a": int64(123)}
	assert.Equal(float64(123), MapGetFloat64(i64, "a", 1))
	i32 := map[string]interface{}{"a": int32(123)}
	assert.Equal(float64(123), MapGetFloat64(i32, "a", 1))

}

func TestMapInsert(t *testing.T) {
	assert := assert.New(t)
	m1 := map[string]interface{}{"a": 1}
	m2 := map[string]interface{}{"a": 2}
	list := []map[string]interface{}{}
	list = MapInsert(list, 0, m1)
	list = MapInsert(list, 0, m2)
	assert.Equal(2, list[0]["a"])
	assert.Equal(1, list[1]["a"])
}
