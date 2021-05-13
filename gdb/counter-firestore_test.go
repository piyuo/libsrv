package gdb

import (
	"context"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/piyuo/libsrv/db"
	"github.com/piyuo/libsrv/identifier"
	"github.com/piyuo/libsrv/test"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	name := "test-counter-" + identifier.RandomString(8)
	counter := client.Counter(name, 1)

	err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		err := counter.IncrementRX(ctx, tx)
		assert.Nil(err)
		return counter.IncrementWX(ctx, tx, 1)
	})
	assert.Nil(err)

	shardsCount, err := counter.ShardsCount(ctx)
	assert.Nil(err)
	assert.Equal(1, shardsCount) // 5 shard, all/year/month/day/hour

	firstCount, err := counter.CountAll(ctx)
	assert.Nil(err)
	assert.True(firstCount > 0)

	//increment again
	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		err := counter.IncrementRX(ctx, tx)
		assert.Nil(err)
		return counter.IncrementWX(ctx, tx, 2)
	})
	assert.Nil(err)

	secondCount, err := counter.CountAll(ctx)
	assert.Nil(err)
	assert.True(secondCount > firstCount)

	err = counter.Delete(ctx)
	assert.Nil(err)

	shardsCount, err = counter.ShardsCount(ctx)
	assert.Nil(err)
	assert.Equal(0, shardsCount)
}

func TestCounterShards(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	name := "test-counter-shards-" + identifier.RandomString(8)
	counter := client.Counter(name, 1)
	defer counter.Delete(ctx)

	err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		err := counter.IncrementRX(ctx, tx)
		assert.Nil(err)
		return counter.IncrementWX(ctx, tx, 1)
	})
	assert.Nil(err)

	utcNow := time.Now().UTC()
	year := strconv.Itoa(utcNow.Year())
	month := strconv.Itoa(int(utcNow.Month()))
	day := strconv.Itoa(int(utcNow.Day()))
	hour := strconv.Itoa(int(utcNow.Hour()))
	c := counter.(*CounterFirestore)
	yearRef := client.(*ClientFirestore).getDocRef(c.collection, c.id+year+"_"+c.pickedShard)
	monthRef := client.(*ClientFirestore).getDocRef(c.collection, c.id+year+"-"+month+"_"+c.pickedShard)
	dayRef := client.(*ClientFirestore).getDocRef(c.collection, c.id+year+"-"+month+"-"+day+"_"+c.pickedShard)
	hourRef := client.(*ClientFirestore).getDocRef(c.collection, c.id+year+"-"+month+"-"+day+"-"+hour+"_"+c.pickedShard)
	hourRef.Delete(ctx)
	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		return counter.IncrementRX(ctx, tx)
	})
	assert.Nil(err)

	dayRef.Delete(ctx)
	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		return counter.IncrementRX(ctx, tx)
	})
	assert.Nil(err)

	monthRef.Delete(ctx)
	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		return counter.IncrementRX(ctx, tx)
	})
	assert.Nil(err)

	yearRef.Delete(ctx)
	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		return counter.IncrementRX(ctx, tx)
	})
	assert.Nil(err)

}

func TestCounterFail(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	name := "test-counter-fail-" + identifier.RandomString(8)
	counter := client.Counter(name, 1)
	defer counter.Delete(ctx)

	firstCount, err := counter.CountAll(ctx)
	assert.Nil(err)

	err = client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		err := counter.IncrementRX(ctx, tx)
		assert.Nil(err)
		err = counter.IncrementWX(ctx, tx, 1)
		assert.Nil(err)
		return errors.New("fail")
	})
	assert.NotNil(err)

	secondCount, err := counter.CountAll(ctx)
	assert.Nil(err)
	assert.Equal(firstCount, secondCount)
}

func TestCounterWxBeforeRx(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	name := "test-counter-wx-before-rx-" + identifier.RandomString(8)
	counter := client.Counter(name, 1)
	defer counter.Delete(ctx)

	err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		return counter.IncrementWX(ctx, tx, 1)
	})
	assert.NotNil(err)
}

func TestCounterInCanceledCtx(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	client := sampleClient()
	name := "test-counter-no-ctx-" + identifier.RandomString(8)
	counter := client.Counter(name, 1)

	ctxCanceled := test.CanceledContext()
	count, err := counter.CountAll(ctxCanceled)
	assert.NotNil(err)
	assert.Equal(float64(0), count)

	err = counter.Delete(ctxCanceled)
	assert.NotNil(err)
}

func TestCounterNumShards0(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	client := sampleClient()
	counter := client.Counter("no-num", 0)
	assert.NotNil(counter)
	assert.Equal(10, counter.(*CounterFirestore).numShards)
}

func TestCounterConcurrent(t *testing.T) {
	t.Parallel()
	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()
	client := sampleClient()
	name := "test-counter-concurrent-" + identifier.RandomString(8)
	counter := client.Counter(name, 30)
	defer counter.Delete(ctx)

	var concurrent = 3
	var wg sync.WaitGroup
	wg.Add(concurrent)
	counting := func() {
		counter := client.Counter(name, 1)
		for i := 0; i < 3; i++ {
			err := client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
				err := counter.IncrementRX(ctx, tx)
				if err != nil {
					t.Errorf("rx err should nil, got %v", err)
					return err
				}
				//fmt.Printf("count:%v\n", i)
				err = counter.IncrementWX(ctx, tx, 1)
				if err != nil {
					t.Errorf("wx err should be nil, got %v", err)
					return err
				}
				return nil
			})
			if err != nil {
				t.Errorf("err should be nil, got %v", err)
				return
			}
		}
		wg.Done()
	}

	//create go routing to do counting
	for i := 0; i < concurrent; i++ {
		go counting()
	}
	wg.Wait()
	count, err := counter.CountAll(ctx)
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
		return
	}
	if count != 9 {
		t.Errorf("count = %v; want 9", count)
	}
}

func TestCounterCountPeriod(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()
	name := "test-counter-count-period-" + identifier.RandomString(8)
	counter := client.Counter(name, 1)
	now := time.Now().UTC()
	from := time.Date(now.Year()-1, 01, 01, 0, 0, 0, 0, time.UTC)
	to := time.Date(now.Year()+1, 01, 01, 0, 0, 0, 0, time.UTC)
	_, err := counter.CountPeriod(ctx, db.HierarchyYear, from, to)
	assert.Contains(err.Error(), "requires an index")
}
