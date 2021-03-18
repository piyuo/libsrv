package gdb

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/src/data"
	"github.com/stretchr/testify/assert"
)

func TestCounters(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()

	counter := g.Counters().Connection.CreateCounter(g.Counters().TableName, "SampleCount", 3, data.DateHierarchyNone)
	defer counter.Clear(ctx)

	err = g.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx)
		assert.Nil(err)
		return counter.IncrementWX(ctx, 1)
	})
	assert.Nil(err)

	count, err := counter.CountAll(ctx)
	assert.Nil(err)
	assert.Equal(float64(1), count)
}
