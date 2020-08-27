package data

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	util "github.com/piyuo/libsrv/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCountPeriod(t *testing.T) {
	Convey("Should init, increment, count on counter", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		countersG, _ := createSampleCounters(dbG, dbR)
		counter := countersG.SampleCounter()
		counterFirestore := counter.(*CounterFirestore)
		err := counter.Clear(ctx)
		So(err, ShouldBeNil)
		defer counter.Clear(ctx)

		// add mock data
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
		So(err, ShouldBeNil)

		from := time.Date(now.Year()-1, 01, 01, 0, 0, 0, 0, time.UTC)
		to := time.Date(now.Year()+1, 01, 01, 0, 0, 0, 0, time.UTC)
		//		fmt.Println(to.Format("2006-01-02 15:04:05"))
		count, err := counter.CountPeriod(ctx, HierarchyYear, from, to)
		So(err, ShouldBeNil)
		So(count, ShouldEqual, 2)
	})
}

func TestCountPeriodR(t *testing.T) {
	Convey("Should init, increment, count on counter", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		_, countersR := createSampleCounters(dbG, dbR)
		counter := countersR.SampleCounter()
		counterFirestore := counter.(*CounterFirestore)
		err := counter.Clear(ctx)
		So(err, ShouldBeNil)
		defer counter.Clear(ctx)

		// add mock data
		now := time.Now().UTC()
		err = dbR.Transaction(ctx, func(ctx context.Context) error {
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
		So(err, ShouldBeNil)

		from := time.Date(now.Year()-1, 01, 01, 0, 0, 0, 0, time.UTC)
		to := time.Date(now.Year()+1, 01, 01, 0, 0, 0, 0, time.UTC)
		//		fmt.Println(to.Format("2006-01-02 15:04:05"))
		count, err := counter.CountPeriod(ctx, HierarchyYear, from, to)
		So(err, ShouldBeNil)
		So(count, ShouldEqual, 2)
	})
}

func TestCounterFailedIncrement(t *testing.T) {
	Convey("Should init, increment, count on counter", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		countersG, _ := createSampleCounters(dbG, dbR)

		counter := countersG.SampleCounter()
		counterFirestore := counter.(*CounterFirestore)

		err := dbG.Transaction(ctx, func(ctx context.Context) error {
			err := counter.IncrementRX(ctx, 1)
			So(err, ShouldBeNil)
			counterFirestore.value = nil // mock error
			return counter.IncrementWX(ctx)
		})
		So(err, ShouldNotBeNil)

		err = dbG.Transaction(ctx, func(ctx context.Context) error {
			err := counter.IncrementRX(ctx, 1)
			So(err, ShouldBeNil)
			counterFirestore.shardExist = true
			counterFirestore.value = nil // mock error
			return counter.IncrementWX(ctx)
		})
		So(err, ShouldNotBeNil)

	})
}

func TestCounterIncrement(t *testing.T) {
	Convey("Should init, increment, count on counter", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		countersG, _ := createSampleCounters(dbG, dbR)

		counter := countersG.SampleCounter()
		So(counter, ShouldNotBeNil)

		err := counter.Clear(ctx)
		So(err, ShouldBeNil)
		defer counter.Clear(ctx)

		err = dbG.Transaction(ctx, func(ctx context.Context) error {
			err := counter.IncrementRX(ctx, 1)
			So(err, ShouldBeNil)
			return counter.IncrementWX(ctx)
		})
		So(err, ShouldBeNil)

		shardsCount, err := counter.ShardsCount(ctx)
		So(err, ShouldBeNil)
		So(shardsCount, ShouldEqual, 1)

		count, err := counter.CountAll(ctx)
		So(err, ShouldBeNil)
		So(count, ShouldEqual, 1)

		//increment again
		err = dbG.Transaction(ctx, func(ctx context.Context) error {
			err := counter.IncrementRX(ctx, 2)
			So(err, ShouldBeNil)
			return counter.IncrementWX(ctx)
		})
		So(err, ShouldBeNil)

		count, err = counter.CountAll(ctx)
		So(err, ShouldBeNil)
		So(count, ShouldEqual, 3)
	})
}

func TestCounter(t *testing.T) {
	Convey("Should init, increment, count on counter", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		countersG, countersR := createSampleCounters(dbG, dbR)

		incrementMustUseWithInTransacton(ctx, dbG, countersG)
		incrementMustUseWithInTransacton(ctx, dbR, countersR)

		testCounter(ctx, dbG, countersG)
		testCounter(ctx, dbR, countersR)

		testCounterClear(ctx, dbG, countersG)
		testCounterClear(ctx, dbR, countersR)

		testCounterInCanceledCtx(ctx, dbG, countersG)
		testCounterInCanceledCtx(ctx, dbR, countersR)
	})
}

func incrementMustUseWithInTransacton(ctx context.Context, db SampleDB, counters *SampleCounters) {
	counter := counters.SampleCounter()
	err := counter.Clear(ctx)
	So(err, ShouldBeNil)
	defer counter.Clear(ctx)

	err = counter.IncrementRX(ctx, 1)
	So(err, ShouldNotBeNil)
	err = counter.IncrementWX(ctx)
	So(err, ShouldNotBeNil)

	num, err := counter.CountAll(ctx)
	So(err, ShouldBeNil)
	So(num, ShouldEqual, 0)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementWX(ctx) // should call Increment First error
		So(err, ShouldNotBeNil)
		return err
	})
	So(err, ShouldNotBeNil)

}

func testCounter(ctx context.Context, db SampleDB, counters *SampleCounters) {
	// create counter
	counter := counters.SampleCounter()

	err := counter.Clear(ctx)
	So(err, ShouldBeNil)
	defer counter.Clear(ctx)

	So(counter, ShouldNotBeNil)
	So((counter.(*CounterFirestore)).errorID(), ShouldNotBeEmpty)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx, 1)
		So(err, ShouldBeNil)
		err = counter.IncrementWX(ctx)
		So(err, ShouldBeNil)
		return err
	})
	So(err, ShouldBeNil)

	num, err := counter.CountAll(ctx)
	So(err, ShouldBeNil)
	So(num, ShouldEqual, 1)

	//try again
	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx, 5)
		So(err, ShouldBeNil)
		return counter.IncrementWX(ctx)
	})
	So(err, ShouldBeNil)

	num, err = counter.CountAll(ctx)
	So(err, ShouldBeNil)
	So(num, ShouldEqual, 6)

	//try minus
	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx, -3)
		So(err, ShouldBeNil)
		return counter.IncrementWX(ctx)
	})
	So(err, ShouldBeNil)

	num, err = counter.CountAll(ctx)
	So(err, ShouldBeNil)
	So(num, ShouldEqual, 3)

	//try count in transaction
	err = db.Transaction(ctx, func(ctx context.Context) error {
		num, err = counter.CountAll(ctx)
		So(err, ShouldBeNil)
		So(num, ShouldEqual, 3)
		return nil
	})

	//counter minimal shards is 10
	counter = counters.Counter("minShards", 0, "UTC", 0)
	So(counter, ShouldNotBeNil)
	firestoreCounter := counter.(*CounterFirestore)
	So(firestoreCounter.numShards, ShouldEqual, 10)
	err = counter.Clear(ctx)
	So(err, ShouldBeNil)
}

func testCounterInCanceledCtx(ctx context.Context, db SampleDB, counters *SampleCounters) {

	counter := counters.SampleCounter()
	So(counter, ShouldNotBeNil)

	err := counter.Clear(ctx)
	So(err, ShouldBeNil)
	defer counter.Clear(ctx)

	ctxCanceled := util.CanceledCtx()

	count, err := counter.CountAll(ctxCanceled)
	So(err, ShouldNotBeNil)
	So(count, ShouldEqual, 0)

	err = counter.Clear(ctxCanceled)
	So(err, ShouldNotBeNil)

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
				err := counter.IncrementRX(ctx, 1)

				if err != nil {
					t.Errorf("err should be nil, got %v", err)
					return err
				}
				//fmt.Printf("count:%v\n", i)
				err = counter.IncrementWX(ctx)
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

func testCounterClear(ctx context.Context, db SampleDB, counters *SampleCounters) {

	counter := counters.SampleCounter()

	err := counter.Clear(ctx)
	So(err, ShouldBeNil)
	defer counter.Clear(ctx)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		counter := counters.SampleCounter()
		err = counter.IncrementRX(ctx, 1)
		So(err, ShouldBeNil)
		return counter.IncrementWX(ctx)
	})
	So(err, ShouldBeNil)

	count, err := counter.CountAll(ctx)
	So(err, ShouldBeNil)
	So(count, ShouldEqual, 1)

	counter.Clear(ctx)

	count, err = counter.CountAll(ctx)
	So(err, ShouldBeNil)
	So(count, ShouldEqual, 0)
}
