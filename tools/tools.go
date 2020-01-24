package tools

import (
	"github.com/lithammer/shortuuid/v3"
)

//UUID generates concise, unambiguous, URL-safe UUIDs
func UUID() string {
	return shortuuid.New()
}
