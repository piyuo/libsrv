package gdb

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/db"
	"github.com/piyuo/libsrv/identifier"
	util "github.com/piyuo/libsrv/util"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()

	name := "test-tx-" + identifier.RandomString(8)
	sample := &Sample{
		Name:  name,
		Value: 1,
	}

	err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		return tx.Set(ctx, sample)
	})
	assert.Nil(err)

	found, err := client.Query(&Sample{}).Where("Name", "==", name).ReturnExists(ctx)
	assert.Nil(err)
	assert.True(found)

	// read before write
	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		sample2, err := tx.Get(ctx, &Sample{}, sample.ID())
		assert.Nil(err)
		assert.NotNil(sample2)

		// test nil safety
		s := sample2.(*Sample)
		assert.NotNil(s.Array)
		assert.NotNil(s.Numbers)
		assert.NotNil(s.PObj)

		found, err := tx.Exists(ctx, &Sample{}, sample.ID())
		assert.Nil(err)
		assert.True(found)

		list, err := tx.List(ctx, &Sample{}, 10)
		assert.Nil(err)
		assert.NotEmpty(list)

		// test nil safety
		s = list[0].(*Sample)
		assert.NotNil(s.Array)
		assert.NotNil(s.Numbers)
		assert.NotNil(s.PObj)

		value, err := tx.Select(ctx, &Sample{}, sample.ID(), "Value")
		assert.Nil(err)
		assert.NotEmpty(list)
		intValue, err := util.ToInt(value)
		assert.Nil(err)
		assert.Equal(1, intValue)
		return nil
	})
	assert.Nil(err)

	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {

		err = tx.Update(ctx, sample, map[string]interface{}{
			"Value": 2,
		})
		assert.Nil(err)
		err = tx.Increment(ctx, sample, "Value", 1)
		assert.Nil(err)
		return nil
	})
	assert.Nil(err)

	// read before write
	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		iSample4, err := tx.Query(&Sample{}).Where("Name", "==", name).ReturnFirst(ctx)
		assert.Nil(err)
		sample4 := iSample4.(*Sample)
		assert.Equal(3, sample4.Value)
		return tx.Delete(ctx, sample)
	})
	assert.Nil(err)
}

func TestTransactionFail(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()

	err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		return errors.New("something wrong")
	})
	assert.NotNil(err)
}

func TestTransactionAssert(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()

	err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		_, err := tx.Get(ctx, nil, "")
		assert.NotNil(err)
		_, err = tx.Get(ctx, &Sample{}, "")
		assert.NotNil(err)
		_, err = tx.Exists(ctx, nil, "")
		assert.NotNil(err)
		_, err = tx.Exists(ctx, &Sample{}, "")
		assert.NotNil(err)
		_, err = tx.List(ctx, nil, 10)
		assert.NotNil(err)
		_, err = tx.Select(ctx, nil, "", "")
		assert.NotNil(err)
		_, err = tx.Select(ctx, &Sample{}, "", "")
		assert.NotNil(err)
		err = tx.Set(ctx, nil)
		assert.NotNil(err)
		err = tx.Update(ctx, nil, nil)
		assert.NotNil(err)
		err = tx.Increment(ctx, nil, "Value", 1)
		assert.NotNil(err)
		err = tx.Delete(ctx, nil)
		assert.NotNil(err)
		_, _, err = tx.(*TransactionFirestore).deleteAll(ctx, nil, 1)
		assert.NotNil(err)

		ref := client.(*ClientFirestore).getDocRef("not-exists", "not-exists")
		err = tx.(*TransactionFirestore).createShard(ref, nil)
		assert.NotNil(err)
		err = tx.(*TransactionFirestore).incrementShard(ref, nil)
		assert.NotNil(err)

		return nil
	})
	assert.Nil(err)
}

func TestTransactionDeleteAll(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()

	name := "test-tx-delete-all-" + identifier.RandomString(8)
	sample := &SampleDeleteAll{
		Name: name,
	}
	err := client.Set(ctx, sample)
	assert.Nil(err)

	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		done, numDeleted, err := tx.(*TransactionFirestore).deleteAll(ctx, sample, 10)
		assert.Nil(err)
		assert.True(done)
		assert.Equal(1, numDeleted)
		return nil
	})
	assert.Nil(err)
}
