package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitTestIsUnitTest(t *testing.T) {
	assert := assert.New(t)
	assert.True(IsUnitTest())
}
