package usage

import (
	"context"
	"testing"
	"time"

	"github.com/piyuo/libsrv/src/data"
	"github.com/stretchr/testify/assert"
)

func NewSample(ctx context.Context) data.DB {
	conn, err := data.FirestoreGlobalConnection(ctx)
	if err != nil {
		panic(err)
	}
	db := &data.BaseDB{
		Connection: conn,
	}
	removeSample(ctx, db)
	return db
}

func removeSample(ctx context.Context, db data.DB) {
	table := &data.Table{
		Connection: db.GetConnection(),
		TableName:  "Usage",
	}
	table.Clear(ctx)
}

func TestUsage(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	db := NewSample(ctx)
	defer removeSample(ctx, db)

	group := "test"
	key := "key1"
	usage := NewUsage(db)

	//check count is 0
	expired := time.Now().UTC().Add(time.Duration(-5) * time.Second)
	count, recent, err := usage.Count(ctx, group, key, expired)
	assert.Nil(err)
	assert.Equal(0, count)
	assert.True(recent.IsZero())

	//add usage
	err = usage.Add(ctx, group, key)
	assert.Nil(err)

	//check count is 1
	count, recent, err = usage.Count(ctx, group, key, expired)
	assert.Nil(err)
	assert.Equal(1, count)
	assert.False(recent.IsZero())
	//		dur := time.Now().UTC().Sub(recent)
	//		sec := dur.Seconds()
	assert.True(time.Now().UTC().After(recent))

	//remove usage
	err = usage.Remove(ctx, group, key)
	assert.Nil(err)

	//check count is 0
	count, recent, err = usage.Count(ctx, group, key, expired)
	assert.Nil(err)
	assert.Equal(0, count)
	assert.True(recent.IsZero())
}

func TestUsageDuration(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	db := NewSample(ctx)
	defer removeSample(ctx, db)

	group := "test"
	key := "key1"
	usage := NewUsage(db)

	//add usage that won't count
	err := usage.Add(ctx, group, key)
	assert.Nil(err)

	time.Sleep(time.Duration(2) * time.Second)

	err = usage.Add(ctx, group, key)
	assert.Nil(err)

	//check count is 1
	expired := time.Now().UTC().Add(time.Duration(-1) * time.Second)
	count, recent, err := usage.Count(ctx, group, key, expired)
	assert.Nil(err)
	assert.Equal(1, count)
	assert.False(recent.IsZero())
	assert.True(time.Now().UTC().After(recent))
	dur := time.Now().UTC().Sub(recent)
	ms := dur.Milliseconds()
	assert.True(ms < 1000)

	//remove usage
	err = usage.Remove(ctx, group, key)
	assert.Nil(err)
}

func TestUsageMaintenance(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	db := NewSample(ctx)
	defer removeSample(ctx, db)

	group := "test"
	key := "key1"
	usage := NewUsage(db)

	//add 2 usage
	err := usage.Add(ctx, group, key)
	assert.Nil(err)

	err = usage.Add(ctx, group, key)
	assert.Nil(err)

	time.Sleep(time.Duration(2) * time.Second)

	// test maintenance usage by remove past 1 seconds usage
	expired := time.Now().UTC().Add(time.Duration(-1) * time.Second)
	result, err := usage.Maintenance(ctx, expired)
	assert.Nil(err)
	assert.True(result)

	count, _, err := usage.Count(ctx, group, key, expired)
	assert.Nil(err)
	assert.Equal(0, count)
}
