package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	timer := NewTimer()
	timer.Start()
	time.Sleep(1 * time.Millisecond)
	ms := timer.TimeSpan()
	assert.True(ms >= 1)
	ms = timer.Stop()
	assert.True(ms >= 1)

}
