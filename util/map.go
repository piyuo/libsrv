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
//	So(m2["a"], ShouldEqual, "1")
//	So(m2["b"], ShouldEqual, "2")
//
func MapFromString(str string) map[string]string {
	ss := strings.Split(str, "&")
	m := make(map[string]string)
	for _, pair := range ss {
		z := strings.Split(pair, "=")
		if len(z) == 1 {
			m[z[0]] = ""

		} else {
			m[z[0]] = z[1]
		}
	}
	return m
}
