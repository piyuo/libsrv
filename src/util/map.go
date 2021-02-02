package util

import (
	"bytes"
	"fmt"
	"strings"
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
