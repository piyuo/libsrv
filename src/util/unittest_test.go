package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitTestInUnitTest(t *testing.T) {
	assert := assert.New(t)
	assert.True(InUnitTest())
}
