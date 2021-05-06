package mapping

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMappingToString(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	m := map[string]interface{}{
		"a": "1",
		"b": "2",
	}
	s := ToString(m)
	assert.NotEmpty(s)

	m2 := FromString(s)
	assert.Equal("1", m2["a"])
	assert.Equal("2", m2["b"])
}

func TestMappingToAndFromString(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	m := map[string]interface{}{}
	s := ToString(m)
	assert.Empty(s)

	m2 := FromString(s)
	assert.Equal(0, len(m2))
}

func TestMappingFromString(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	m := FromString("=")
	assert.Equal(0, len(m))
}

func TestMappingGetMap(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	m := map[string]interface{}{"a": map[string]interface{}{"c": "d"}, "b": 1}
	a := GetMap(m, "a")
	assert.Equal("d", GetString(a, "c", ""))
	//map not exists
	aNil := GetMap(m, "x")
	assert.Empty(aNil)
	//wrong data type
	b := GetMap(m, "b")
	assert.Empty(b)
}

func TestMappingGetString(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	m := map[string]interface{}{"a": "b", "c": 1}
	assert.Equal("b", GetString(m, "a", "default"))
	assert.Equal("default", GetString(m, "b", "default"))
	assert.Equal("1", GetString(m, "c", "default"))
}

func TestMappingGetFloat64(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	m := map[string]interface{}{"a": 123, "c": "b"}
	assert.Equal(float64(123), GetFloat64(m, "a", 1))
	assert.Equal(float64(0), GetFloat64(m, "b", 0))
	assert.Equal(float64(0), GetFloat64(m, "c", 0))

	f64 := map[string]interface{}{"a": float64(123)}
	assert.Equal(float64(123), GetFloat64(f64, "a", 1))
	f32 := map[string]interface{}{"a": float32(123)}
	assert.Equal(float64(123), GetFloat64(f32, "a", 1))
	i64 := map[string]interface{}{"a": int64(123)}
	assert.Equal(float64(123), GetFloat64(i64, "a", 1))
	i32 := map[string]interface{}{"a": int32(123)}
	assert.Equal(float64(123), GetFloat64(i32, "a", 1))

}

func TestMappingInsert(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	m1 := map[string]interface{}{"a": 1}
	m2 := map[string]interface{}{"a": 2}
	list := []interface{}{}
	list = Insert(list, 0, m1)
	list = Insert(list, 0, m2)
	assert.Equal(2, list[0].(map[string]interface{})["a"])
	assert.Equal(1, list[1].(map[string]interface{})["a"])
}

func TestMappingGetTime(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	j := map[string]interface{}{
		"a": time.Date(2020, time.April, 11, 21, 34, 01, 0, time.UTC)}

	//test return time directly
	d := GetTime(j, "a", time.Time{})
	assert.Equal(2020, d.Year())
	assert.Equal(11, d.Day())

	// test time after marshal
	bytes, err := json.Marshal(j)
	assert.Nil(err)
	err = json.Unmarshal(bytes, &j)
	assert.Nil(err)

	d = GetTime(j, "a", time.Time{})
	assert.Equal(2020, d.Year())
	assert.Equal(11, d.Day())

	//test not exist
	d = GetTime(j, "b", time.Time{})
	assert.True(d.IsZero())

	//test time string in wrong format
	j["a"] = "not-time-format"
	d = GetTime(j, "a", time.Time{})
	assert.True(d.IsZero())

	//test time string in wrong data type
	j["a"] = 123
	d = GetTime(j, "a", time.Time{})
	assert.True(d.IsZero())

}

func TestMappingGetList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	j := map[string]interface{}{
		"a": []interface{}{"b", "c"},
	}
	list := j["a"].([]interface{})
	assert.Len(list, 2)
}
