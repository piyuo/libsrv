package data

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	util "github.com/piyuo/libsrv/util"
	"github.com/stretchr/testify/assert"
)

func TestCountPeriod(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()

	counter := g.Counters().SampleCounter()
	defer counter.Clear(ctx)
	counterFirestore := counter.(*CounterFirestore)
	assert.Nil(err)

	// mock data
	now := time.Now().UTC()
	err = g.Transaction(ctx, func(ctx context.Context) error {
		if err := counterFirestore.mock(HierarchyYear, now, 1, 1); err != nil {
			return err
		}
		if err := counterFirestore.mock(HierarchyYear, now.AddDate(-1, 0, 0), 2, 1); err != nil {
			return err
		}
		if err := counterFirestore.mock(HierarchyYear, now.AddDate(-2, 0, 0), 3, 1); err != nil {
			return err
		}
		return nil
	})
	assert.Nil(err)

	from := time.Date(now.Year()-1, 01, 01, 0, 0, 0, 0, time.UTC)
	to := time.Date(now.Year()+1, 01, 01, 0, 0, 0, 0, time.UTC)
	count, err := counter.CountPeriod(ctx, HierarchyYear, from, to)
	assert.Nil(err)
	assert.Equal(float64(2), count)

	//test DetailPeriod
	dict, err := counter.DetailPeriod(ctx, HierarchyYear, from, to)
	assert.Equal(2, len(dict))
}

func TestCounterFailedIncrement(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()

	counter := g.Counters().SampleCounter()
	defer counter.Clear(ctx)
	counterFirestore := counter.(*CounterFirestore)
	assert.Nil(err)

	err = g.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx)
		assert.Nil(err)
		counterFirestore.callRX = false // mock error
		return counter.IncrementWX(ctx, 1)
	})
	assert.NotNil(err)

	err = g.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx)
		assert.Nil(err)
		counterFirestore.shardAllExist = true
		counterFirestore.callRX = false // mock error
		return counter.IncrementWX(ctx, 1)
	})
	assert.NotNil(err)
}

func TestCounterIncrement(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()

	counter := g.Counters().Counter("SampleCount", 3, DateHierarchyFull)
	defer counter.Clear(ctx)
	assert.NotNil(counter)

	err = g.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx)
		assert.Nil(err)
		return counter.IncrementWX(ctx, 1)
	})
	assert.Nil(err)

	shardsCount, err := counter.ShardsCount(ctx)
	assert.Nil(err)
	assert.Equal(5, shardsCount) // 5 shard, all/year/month/day/hour

	count, err := counter.CountAll(ctx)
	assert.Nil(err)
	assert.Equal(float64(1), count)

	//increment again
	err = g.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx)
		assert.Nil(err)
		return counter.IncrementWX(ctx, 2)
	})
	assert.Nil(err)

	count, err = counter.CountAll(ctx)
	assert.Nil(err)
	assert.Equal(float64(3), count)
}

func TestIncrementMustUseWithInTransacton(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()

	counter := g.Counters().SampleCounter()
	defer counter.Clear(ctx)

	err = counter.IncrementRX(ctx)
	assert.NotNil(err)
	err = counter.IncrementWX(ctx, 1)
	assert.NotNil(err)

	num, err := counter.CountAll(ctx)
	assert.Nil(err)
	assert.Equal(float64(0), num)

	err = g.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementWX(ctx, 1) // should call Increment First error
		assert.NotNil(err)
		return err
	})
	assert.NotNil(err)
}

func TestCounter(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()

	counter := g.Counters().SampleCounter()
	defer counter.Clear(ctx)

	assert.NotEmpty((counter.(*CounterFirestore)).errorID())

	err = g.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx)
		assert.Nil(err)
		err = counter.IncrementWX(ctx, 1)
		assert.Nil(err)
		return err
	})
	assert.Nil(err)

	num, err := counter.CountAll(ctx)
	assert.Nil(err)
	assert.Equal(float64(1), num)

	//try again
	err = g.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx)
		assert.Nil(err)
		return counter.IncrementWX(ctx, 5)
	})
	assert.Nil(err)

	num, err = counter.CountAll(ctx)
	assert.Nil(err)
	assert.Equal(float64(6), num)

	//try minus
	err = g.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx)
		assert.Nil(err)
		return counter.IncrementWX(ctx, -3)
	})
	assert.Nil(err)

	num, err = counter.CountAll(ctx)
	assert.Nil(err)
	assert.Equal(float64(3), num)

	//try count in transaction
	err = g.Transaction(ctx, func(ctx context.Context) error {
		num, err = counter.CountAll(ctx)
		assert.Nil(err)
		assert.Equal(float64(3), num)
		return nil
	})

	//counter minimal shards is 10
	counter = g.Counters().Counter("minShards", 0, DateHierarchyNone)
	assert.NotNil(counter)
	firestoreCounter := counter.(*CounterFirestore)
	assert.Equal(10, firestoreCounter.numShards)
	err = counter.Clear(ctx)
	assert.Nil(err)
}

func TestCounterInCanceledCtx(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()

	counter := g.Counters().SampleCounter()
	defer counter.Clear(ctx)

	ctxCanceled := util.CanceledCtx()

	count, err := counter.CountAll(ctxCanceled)
	assert.NotNil(err)
	assert.Equal(float64(0), count)

	err = counter.Clear(ctxCanceled)
	assert.NotNil(err)
}

func TestConcurrentCounter(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()
	dbG, _ := NewSampleGlobalDB(ctx)
	defer dbG.Close()
	countersG := dbG.Counters()

	counter := countersG.SampleCounter()
	defer counter.Clear(ctx)

	var concurrent = 3
	var wg sync.WaitGroup
	wg.Add(concurrent)
	counting := func() {
		db, _ := NewSampleGlobalDB(ctx)
		defer db.Close()
		counter := db.Counters().SampleCounter()
		for i := 0; i < 3; i++ {
			err := db.Transaction(ctx, func(ctx context.Context) error {
				err := counter.IncrementRX(ctx)

				if err != nil {
					t.Errorf("err should be nil, got %v", err)
					return err
				}
				//fmt.Printf("count:%v\n", i)
				err = counter.IncrementWX(ctx, 1)
				if err != nil {
					t.Errorf("err should be nil, got %v", err)
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

func TestCounterClear(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	defer g.Close()

	counter := g.Counters().SampleCounter()
	defer counter.Clear(ctx)

	err = g.Transaction(ctx, func(ctx context.Context) error {
		c := g.Counters().SampleCounter()
		err = c.IncrementRX(ctx)
		assert.Nil(err)
		return c.IncrementWX(ctx, 1)
	})
	assert.Nil(err)

	count, err := counter.CountAll(ctx)
	assert.Nil(err)
	assert.Equal(float64(1), count)

	counter.Clear(ctx)

	count, err = counter.CountAll(ctx)
	assert.Nil(err)
	assert.Equal(float64(0), count)
}
