package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanceledContext(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := CanceledContext()
	assert.NotNil(ctx.Err())
}

/*
var tests = []struct {
		input    int
		expected int
	}{
		{2, 4},
		{-1, 1},
		{0, 2},
		{-5, -3},
		{99999, 100001},
	}

	for _, test := range tests {
		assert.Equal(Calculate(test.input), test.expected)
	}
*/
