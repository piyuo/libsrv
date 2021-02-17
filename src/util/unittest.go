package util

import (
	"os"
	"strings"
)

// InUnitTest return true if is in go unit test
//
func InUnitTest() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}
