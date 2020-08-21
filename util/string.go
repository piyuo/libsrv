package util

import (
	"hash/fnv"
	"strings"
)

// StringBetween Get substring between two strings
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
func StringBefore(value string, a string) string {
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}

// StringAfter Get substring after a string
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
	return value[adjustedPos:len(value)]
}

// ArrayToString convert list of string to string
//
//	array := []string{"1", "2", "3"}
//	str := ArrayToString(array) //1,2,3
//	ary := StringToArray(str)
//
func ArrayToString(stringArray []string) string {
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
func StringHash(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return h.Sum32()
}
