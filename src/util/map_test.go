package util

import (
	"encoding/json"
	"testing"
	"time"

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

func TestMapGetMap(t *testing.T) {
	assert := assert.New(t)
	m := map[string]interface{}{"a": map[string]interface{}{"c": "d"}, "b": 1}
	a := MapGetMap(m, "a")
	assert.Equal("d", MapGetString(a, "c", ""))
	//map not exists
	aNil := MapGetMap(m, "x")
	assert.Empty(aNil)
	//wrong data type
	b := MapGetMap(m, "b")
	assert.Empty(b)
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
	list := []interface{}{}
	list = MapInsert(list, 0, m1)
	list = MapInsert(list, 0, m2)
	assert.Equal(2, list[0].(map[string]interface{})["a"])
	assert.Equal(1, list[1].(map[string]interface{})["a"])
}

func TestMapGetTime(t *testing.T) {
	assert := assert.New(t)
	j := map[string]interface{}{
		"a": time.Date(2020, time.April, 11, 21, 34, 01, 0, time.UTC)}

	//test return time directly
	d := MapGetTime(j, "a", time.Time{})
	assert.Equal(2020, d.Year())
	assert.Equal(11, d.Day())

	// test time after marshal
	bytes, err := json.Marshal(j)
	assert.Nil(err)
	err = json.Unmarshal(bytes, &j)
	assert.Nil(err)

	d = MapGetTime(j, "a", time.Time{})
	assert.Equal(2020, d.Year())
	assert.Equal(11, d.Day())

	//test not exist
	d = MapGetTime(j, "b", time.Time{})
	assert.True(d.IsZero())

	//test time string in wrong format
	j["a"] = "not-time-format"
	d = MapGetTime(j, "a", time.Time{})
	assert.True(d.IsZero())

	//test time string in wrong data type
	j["a"] = 123
	d = MapGetTime(j, "a", time.Time{})
	assert.True(d.IsZero())

}

func TestMapGetList(t *testing.T) {
	assert := assert.New(t)
	j := map[string]interface{}{
		"a": []interface{}{"b", "c"},
	}
	list := j["a"].([]interface{})
	assert.Len(list, 2)
}
