package data

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounters(t *testing.T) {
	dbG, dbR := createSampleDB()
	defer removeSampleDB(dbG, dbR)
	cg, _ := createSampleCounters(dbG, dbR)

	countersTest(dbG, t, cg)
}

func countersTest(db SampleDB, t *testing.T, counters *SampleCounters) {
	assert := assert.New(t)
	ctx := context.Background()

	counter := counters.Counter("SampleCount", 3, DateHierarchyNone)
	assert.NotNil(counter)

	err := counter.Clear(ctx)
	assert.Nil(err)
	defer counter.Clear(ctx)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx)
		assert.Nil(err)
		return counter.IncrementWX(ctx, 1)
	})
	assert.Nil(err)

	count, err := counter.CountAll(ctx)
	assert.Nil(err)
	assert.Equal(float64(1), count)
}
