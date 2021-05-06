package mapping

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// ToString convert map to string using key value pair
//
//	m := map[string]interface{}{
//		"a": "1",
//		"b": "2",
//	}
//	s := ToString(m)
//
//
func ToString(m map[string]interface{}) string {
	if len(m) == 0 {
		return ""
	}

	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=%v&", key, value)
	}
	str := b.String()
	str = str[0 : len(str)-1]
	return str
}

// FromString convert key value pair string to map
//
//	m2 := FromString(s)
//
func FromString(str string) map[string]interface{} {
	ss := strings.Split(str, "&")
	m := make(map[string]interface{})
	for _, pair := range ss {
		if pair != "" {
			z := strings.Split(pair, "=")
			if len(z) == 2 && z[0] != "" && z[1] != "" {
				m[z[0]] = z[1]
			}
		}
	}
	return m
}

// GetMap get string value from map, return default value if not exist
//
//	m := GetMap(Map,"submap")
//
func GetMap(mapping map[string]interface{}, key string) map[string]interface{} {
	value := mapping[key]
	if value == nil {
		return map[string]interface{}{}
	}
	switch value.(type) {
	case map[string]interface{}:
		return value.(map[string]interface{})
	}
	return map[string]interface{}{}
}

// GetString get string value from map, return default value if not exist
//
//	str := GetString(Map,"hello")
//
func GetString(mapping map[string]interface{}, key, defaultValue string) string {
	value := mapping[key]
	if value == nil {
		return defaultValue
	}
	return fmt.Sprint(value)
}

// GetFloat64 get float64 value from map, return default value if not exist
//
//	str := GetFloat64(Map,0)
//
func GetFloat64(mapping map[string]interface{}, key string, defaultValue float64) float64 {
	value := mapping[key]
	if value == nil {
		return defaultValue
	}
	switch i := value.(type) {
	case float64:
		return i
	case float32:
		return float64(i)
	case int:
		return float64(i)
	case int32:
		return float64(i)
	case int64:
		return float64(i)
	}
	return defaultValue
}

// GetTime get time value from map, return default value if not exist or wrong time type
//
//	str := MapGetFloat64(Map,time.Now())
//
func GetTime(mapping map[string]interface{}, key string, defaultValue time.Time) time.Time {
	value := mapping[key]
	if value == nil {
		return defaultValue
	}
	switch i := value.(type) {
	case string:
		layout := "2006-01-02T15:04:05Z"
		t, err := time.Parse(layout, string(i))
		if err != nil {
			return defaultValue
		}
		return t
	case time.Time:
		return time.Time(i)
	}
	return defaultValue
}

// Insert insert json object to array
//
func Insert(a []interface{}, index int, value map[string]interface{}) []interface{} {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}
