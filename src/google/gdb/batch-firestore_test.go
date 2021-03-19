package gdb

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/src/db"
	"github.com/piyuo/libsrv/src/identifier"
	"github.com/stretchr/testify/assert"
)

func TestGdbBatch(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()

	name := "testGdb-batch" + identifier.RandomString(6)
	sample := &Sample{
		Name:  name,
		Value: 2001,
	}

	err := client.Batch(ctx, func(ctx context.Context, batch db.Batch) error {
		batch.Set(ctx, sample) //batch mode do not return error
		return nil
	})
	assert.Nil(err)

	count, err := client.Query(&Sample{}).ReturnCount(ctx)
	assert.Nil(err)
	assert.True(count >= 1)
	assert.NotEmpty(sample.ID)

	err = client.Batch(ctx, func(ctx context.Context, batch db.Batch) error {
		batch.Update(ctx, sample, map[string]interface{}{
			"Value": 2002,
		})
		return nil
	})
	assert.Nil(err)

	iSample1, err := client.Query(&Sample{}).Where("Name", "==", name).ReturnFirst(ctx)
	sample1 := iSample1.(*Sample)
	assert.Nil(err)
	assert.Equal(sample1.Value, 2002)

	err = client.Batch(ctx, func(ctx context.Context, batch db.Batch) error {
		batch.Increment(ctx, sample1, "Value", 1)
		return nil
	})
	assert.Nil(err)

	iSample2, err := client.Query(&Sample{}).Where("Name", "==", name).ReturnFirst(ctx)
	sample2 := iSample2.(*Sample)
	assert.Nil(err)
	assert.Equal(sample2.Value, 2003)

	err = client.Batch(ctx, func(ctx context.Context, batch db.Batch) error {
		batch.Delete(ctx, sample2)
		return nil
	})
	assert.Nil(err)

	found, err := client.Query(&Sample{}).Where("Name", "==", name).ReturnIsExists(ctx)
	assert.Nil(err)
	assert.False(found)
}

func TestGdbBatchEmpty(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	// do nothing batch should not result error
	err := client.Batch(ctx, func(ctx context.Context, batch db.Batch) error {
		return nil
	})
	assert.Nil(err)
}

func TestGdbBatchDeleteList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()

	name := "testGdb-batch" + identifier.RandomString(6)
	sample := &Sample{
		Name:  name,
		Value: 2010,
	}
	err := client.Set(ctx, sample)
	assert.Nil(err)

	err = client.Batch(ctx, func(ctx context.Context, batch db.Batch) error {
		batch.DeleteList(ctx, &Sample{}, []string{sample.ID()}) //batch mode do not return error
		return nil
	})

	found, err := client.Query(&Sample{}).Where("Name", "==", name).ReturnIsExists(ctx)
	assert.Nil(err)
	assert.False(found)
}
