package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitTestIsUnitTest(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	assert.True(IsUnitTest())
}
