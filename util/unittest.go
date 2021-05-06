package util

import (
	"os"
	"strings"
)

// IsUnitTest return true if is in go unit test
//
func IsUnitTest() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}
