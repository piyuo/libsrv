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
	dbG, dbR := createSampleDB()
	defer removeSampleDB(dbG, dbR)
	countersG, _ := createSampleCounters(dbG, dbR)
	counter := countersG.SampleCounter()
	counterFirestore := counter.(*CounterFirestore)
	err := counter.Clear(ctx)
	assert.Nil(err)
	defer counter.Clear(ctx)

	// mock data
	now := time.Now().UTC()
	err = dbG.Transaction(ctx, func(ctx context.Context) error {
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
	dbG, dbR := createSampleDB()
	defer removeSampleDB(dbG, dbR)
	countersG, _ := createSampleCounters(dbG, dbR)

	counter := countersG.SampleCounter()
	counterFirestore := counter.(*CounterFirestore)

	err := dbG.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx)
		assert.Nil(err)
		counterFirestore.callRX = false // mock error
		return counter.IncrementWX(ctx, 1)
	})
	assert.NotNil(err)

	err = dbG.Transaction(ctx, func(ctx context.Context) error {
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
	dbG, dbR := createSampleDB()
	defer removeSampleDB(dbG, dbR)
	countersG, _ := createSampleCounters(dbG, dbR)

	counter := countersG.Counter("SampleCount", 3, DateHierarchyFull)
	assert.NotNil(counter)

	err := counter.Clear(ctx)
	assert.Nil(err)
	defer counter.Clear(ctx)

	err = dbG.Transaction(ctx, func(ctx context.Context) error {
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
	err = dbG.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx)
		assert.Nil(err)
		return counter.IncrementWX(ctx, 2)
	})
	assert.Nil(err)

	count, err = counter.CountAll(ctx)
	assert.Nil(err)
	assert.Equal(float64(3), count)
}

func TestCounter(t *testing.T) {
	ctx := context.Background()
	dbG, dbR := createSampleDB()
	defer removeSampleDB(dbG, dbR)
	countersG, _ := createSampleCounters(dbG, dbR)

	incrementMustUseWithInTransacton(ctx, t, dbG, countersG)
	testCounter(ctx, t, dbG, countersG)
	testCounterClear(ctx, t, dbG, countersG)
	testCounterInCanceledCtx(ctx, t, dbG, countersG)
}

func incrementMustUseWithInTransacton(ctx context.Context, t *testing.T, db SampleDB, counters *SampleCounters) {
	assert := assert.New(t)
	counter := counters.SampleCounter()
	err := counter.Clear(ctx)
	assert.Nil(err)
	defer counter.Clear(ctx)

	err = counter.IncrementRX(ctx)
	assert.NotNil(err)
	err = counter.IncrementWX(ctx, 1)
	assert.NotNil(err)

	num, err := counter.CountAll(ctx)
	assert.Nil(err)
	assert.Equal(float64(0), num)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementWX(ctx, 1) // should call Increment First error
		assert.NotNil(err)
		return err
	})
	assert.NotNil(err)
}

func testCounter(ctx context.Context, t *testing.T, db SampleDB, counters *SampleCounters) {
	assert := assert.New(t)
	// create counter
	counter := counters.SampleCounter()

	err := counter.Clear(ctx)
	assert.Nil(err)
	defer counter.Clear(ctx)

	assert.NotNil(counter)
	assert.NotEmpty((counter.(*CounterFirestore)).errorID())

	err = db.Transaction(ctx, func(ctx context.Context) error {
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
	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx)
		assert.Nil(err)
		return counter.IncrementWX(ctx, 5)
	})
	assert.Nil(err)

	num, err = counter.CountAll(ctx)
	assert.Nil(err)
	assert.Equal(float64(6), num)

	//try minus
	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx)
		assert.Nil(err)
		return counter.IncrementWX(ctx, -3)
	})
	assert.Nil(err)

	num, err = counter.CountAll(ctx)
	assert.Nil(err)
	assert.Equal(float64(3), num)

	//try count in transaction
	err = db.Transaction(ctx, func(ctx context.Context) error {
		num, err = counter.CountAll(ctx)
		assert.Nil(err)
		assert.Equal(float64(3), num)
		return nil
	})

	//counter minimal shards is 10
	counter = counters.Counter("minShards", 0, DateHierarchyNone)
	assert.NotNil(counter)
	firestoreCounter := counter.(*CounterFirestore)
	assert.Equal(10, firestoreCounter.numShards)
	err = counter.Clear(ctx)
	assert.Nil(err)
}

func testCounterInCanceledCtx(ctx context.Context, t *testing.T, db SampleDB, counters *SampleCounters) {
	assert := assert.New(t)
	counter := counters.SampleCounter()
	assert.NotNil(counter)

	err := counter.Clear(ctx)
	assert.Nil(err)
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
	err := counter.Clear(ctx)
	defer counter.Clear(ctx)

	var concurrent = 3
	var wg sync.WaitGroup
	wg.Add(concurrent)
	counting := func() {
		db, _ := NewSampleGlobalDB(ctx)
		defer db.Close()
		counter := db.Counters().SampleCounter()
		for i := 0; i < 5; i++ {
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
	if count != 15 {
		t.Errorf("count = %v; want 15", count)
	}
}

func testCounterClear(ctx context.Context, t *testing.T, db SampleDB, counters *SampleCounters) {
	assert := assert.New(t)

	counter := counters.SampleCounter()
	err := counter.Clear(ctx)
	assert.Nil(err)
	defer counter.Clear(ctx)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		counter := counters.SampleCounter()
		err = counter.IncrementRX(ctx)
		assert.Nil(err)
		return counter.IncrementWX(ctx, 1)
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
