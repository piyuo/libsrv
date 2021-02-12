package util

import (
	"hash/fnv"
	"strings"
)

// StringBetween Get substring between two strings
//
//	assert.Equal("2", StringBetween("123", "1", "3"))
//
func StringBetween(value string, a string, b string) string {
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

// StringBefore Get substring before a string
//
//	assert.Equal("1", StringBefore("123", "2"))
//
func StringBefore(value string, a string) string {
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}

// StringAfter Get substring after a string
//
//	assert.Equal("3", StringAfter("123", "2"))
//
func StringAfter(value string, a string) string {
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:]
}

// StringFromArray convert list of string to string
//
//	array := []string{"1", "2", "3"}
//	str := StringFromArray(array) //1,2,3
//	ary := StringToArray(str)
//
func StringFromArray(stringArray []string) string {
	return strings.Join(stringArray, ",")
}

// StringToArray split string to []string
//
//	array := []string{"1", "2", "3"}
//	str := ArrayToString(array)  //1,2,3
//	ary := StringToArray(str)
//
func StringToArray(str string) []string {
	return strings.Split(str, ",")
}

// StringHash Get hash code for string
//
//	code := StringHash(str)
//
func StringHash(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return h.Sum32()
}

// StringsRemove remove value from string array
//
//	ary := []string{"a", "", "b"}
//	filtered := StringsRemove(ary,"")
//
func StringsRemove(s []string, value string) []string {
	var r []string
	for _, str := range s {
		if str != value {
			r = append(r, str)
		}
	}
	return r
}

// StringsContain takes a slice and looks for an element in it, return true if exist otherwise it will return false
//
//	exist :=StringsContain(ary, "a")
//
func StringsContain(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
