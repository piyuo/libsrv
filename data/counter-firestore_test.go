package data

import (
	"context"
	"fmt"
	"sync"
	"testing"

	util "github.com/piyuo/libsrv/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCounter(t *testing.T) {
	Convey("Should init, increment, count on counter", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		countersG, countersR := createSampleCounters(dbG, dbR)
		defer removeSampleCounters(countersG, countersR)

		incrementMustUseWithInTransacton(ctx, dbG, countersG)
		incrementMustUseWithInTransacton(ctx, dbR, countersR)

		testCounter(ctx, dbG, countersG)
		testCounter(ctx, dbR, countersR)

		testCounterReset(ctx, dbG, countersG)
		testCounterReset(ctx, dbR, countersR)

		testCounterInCanceledCtx(ctx, dbG, countersG)
		testCounterInCanceledCtx(ctx, dbR, countersR)
	})
}

func incrementMustUseWithInTransacton(ctx context.Context, db SampleDB, counters *SampleCounters) {
	counter := counters.SampleCounter()
	// clean counter
	err := counters.DeleteSampleCounter(ctx)
	So(err, ShouldBeNil)

	err = counter.IncrementRX(ctx, 1)
	So(err, ShouldNotBeNil)
	err = counter.IncrementWX(ctx)
	So(err, ShouldNotBeNil)

	num, err := counter.Count(ctx)
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
	// clean counter
	err := counters.DeleteSampleCounter(ctx)
	So(err, ShouldBeNil)

	// create counter
	counter := counters.SampleCounter()
	So(counter, ShouldNotBeNil)
	So((counter.(*CounterFirestore)).errorID(), ShouldNotBeEmpty)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx, 1)
		So(err, ShouldBeNil)
		return counter.IncrementWX(ctx)
	})
	So(err, ShouldBeNil)

	num, err := counter.Count(ctx)
	So(err, ShouldBeNil)
	So(num, ShouldEqual, 1)

	//try again
	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx, 5)
		So(err, ShouldBeNil)
		return counter.IncrementWX(ctx)
	})
	So(err, ShouldBeNil)

	num, err = counter.Count(ctx)
	So(err, ShouldBeNil)
	So(num, ShouldEqual, 6)

	//try minus
	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx, -3)
		So(err, ShouldBeNil)
		return counter.IncrementWX(ctx)
	})
	So(err, ShouldBeNil)

	num, err = counter.Count(ctx)
	So(err, ShouldBeNil)
	So(num, ShouldEqual, 3)

	//try count in transaction
	err = db.Transaction(ctx, func(ctx context.Context) error {
		num, err = counter.Count(ctx)
		So(err, ShouldBeNil)
		So(num, ShouldEqual, 3)
		return nil
	})

	//counter minimal shards is 10
	counter = counters.Counter("minShards", 0)
	So(counter, ShouldNotBeNil)
	firestoreCounter := counter.(*CounterFirestore)
	So(firestoreCounter.numShards, ShouldEqual, 10)
	err = counters.Delete(ctx, "minShards")
	So(err, ShouldBeNil)

	// delete exist counter
	err = counters.DeleteSampleCounter(ctx)
	So(err, ShouldBeNil)

}

func testCounterInCanceledCtx(ctx context.Context, db SampleDB, counters *SampleCounters) {

	counter := counters.SampleCounter()
	So(counter, ShouldNotBeNil)

	ctxCanceled := util.CanceledCtx()

	count, err := counter.Count(ctxCanceled)
	So(err, ShouldNotBeNil)
	So(count, ShouldEqual, 0)

	err = counter.Reset(ctxCanceled)
	So(err, ShouldNotBeNil)

	err = counters.DeleteSampleCounter(ctxCanceled)
	So(err, ShouldNotBeNil)
}

func TestConcurrentCounter(t *testing.T) {
	ctx := context.Background()
	dbG, _ := NewSampleGlobalDB(ctx)
	defer dbG.Close()
	countersG := dbG.Counters()
	countersG.DeleteSampleCounter(ctx)

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
				fmt.Printf("count:%v\n", i)
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
	counter := countersG.SampleCounter()
	count, err := counter.Count(ctx)
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
		return
	}
	if count != 15 {
		t.Errorf("count = %v; want 15", count)
	}
}

func testCounterReset(ctx context.Context, db SampleDB, counters *SampleCounters) {
	err := counters.DeleteSampleCounter(ctx)
	defer counters.DeleteSampleCounter(ctx)
	So(err, ShouldBeNil)

	iCounter := counters.SampleCounter()
	ctr := iCounter.(*CounterFirestore)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		counter := counters.SampleCounter()
		err = counter.IncrementRX(ctx, 1)
		So(err, ShouldBeNil)
		return counter.IncrementWX(ctx)
	})
	So(err, ShouldBeNil)

	docCount, shardsCount, err := ctr.shardsInfo(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 1)
	So(shardsCount, ShouldEqual, 1)

	// reset
	counter := counters.SampleCounter()
	counter.Reset(ctx)

	docCount, shardsCount, err = ctr.shardsInfo(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 0)
	So(shardsCount, ShouldEqual, 0)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		counter := counters.SampleCounter()
		err = counter.IncrementRX(ctx, 1)
		So(err, ShouldBeNil)
		return counter.IncrementWX(ctx)
	})
	So(err, ShouldBeNil)

	docCount, shardsCount, err = ctr.shardsInfo(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 1)
	So(shardsCount, ShouldEqual, 1)

	// reset in transaction
	counter = counters.SampleCounter()
	err = db.Transaction(ctx, func(ctx context.Context) error {
		return counter.Reset(ctx)
	})

	docCount, shardsCount, err = ctr.shardsInfo(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 0)
	So(shardsCount, ShouldEqual, 0)
}
