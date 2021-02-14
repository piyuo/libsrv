package util

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// MapToString convert map to string using key value pair
//
//	m := map[string]string{
//		"a": "1",
//		"b": "2",
//	}
//	s := MapToString(m)
//
//
func MapToString(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}

	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=%s&", key, value)
	}
	str := b.String()
	str = str[0 : len(str)-1]
	return str
}

// MapFromString convert key value pair string to map
//
//	m2 := MapFromString(s)
//
func MapFromString(str string) map[string]string {
	ss := strings.Split(str, "&")
	m := make(map[string]string)
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

// MapGetMap get string value from map, return default value if not exist
//
//	m := MapGetMap(Map,"submap")
//
func MapGetMap(mapping map[string]interface{}, key string) map[string]interface{} {
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

// MapGetString get string value from map, return default value if not exist
//
//	str := MapGetString(Map,"hello")
//
func MapGetString(mapping map[string]interface{}, key, defaultValue string) string {
	value := mapping[key]
	if value == nil {
		return defaultValue
	}
	return fmt.Sprint(value)
}

// MapGetFloat64 get float64 value from map, return default value if not exist
//
//	str := MapGetFloat64(Map,0)
//
func MapGetFloat64(mapping map[string]interface{}, key string, defaultValue float64) float64 {
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

// MapGetTime get time value from map, return default value if not exist or wrong time type
//
//	str := MapGetFloat64(Map,time.Now())
//
func MapGetTime(mapping map[string]interface{}, key string, defaultValue time.Time) time.Time {
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

// MapInsert insert json object to array
//
func MapInsert(a []interface{}, index int, value map[string]interface{}) []interface{} {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}
