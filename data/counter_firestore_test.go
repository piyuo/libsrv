package data

import (
	"context"
	"sync"
	"testing"

	util "github.com/piyuo/libsrv/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCounter(t *testing.T) {
	Convey("Should init, increment, count on counter", t, func() {
		ctx := context.Background()
		dbG, dbR, samplesG, samplesR := firestoreBeginTest()
		defer dbG.Close()
		defer dbR.Close()
		countersG := dbG.Counters()
		countersR := dbR.Counters()

		testCounterWithoutCreateShards(ctx, dbG, countersG)
		testCounterWithoutCreateShards(ctx, dbR, countersR)

		testCounter(ctx, dbG, countersG)
		testCounter(ctx, dbR, countersR)

		testCounterInCanceledCtx(ctx, dbG, countersG)
		testCounterInCanceledCtx(ctx, dbR, countersR)

		testCounterInTransaction(ctx, dbG, countersG)
		testCounterInTransaction(ctx, dbR, countersR)

		firestoreEndTest(dbG, dbR, samplesG, samplesR)
	})
}

func testCounterWithoutCreateShards(ctx context.Context, db SampleDB, counters *SampleCounters) {
	// clean counter
	err := counters.DeleteSampleTotal(ctx)
	So(err, ShouldBeNil)

	// test create all shards
	counter := counters.SampleTotal(ctx)
	count, err := counter.Count(ctx)
	So(err, ShouldNotBeNil)
	So(count, ShouldEqual, 0)
	err = counter.Increment(ctx, 1)
	So(err, ShouldNotBeNil)

	// test create all shards in transaction
	err = db.Transaction(ctx, func(ctx context.Context) error {
		counter := counters.SampleTotal(ctx)
		count, err := counter.Count(ctx)
		So(err, ShouldNotBeNil)
		So(count, ShouldEqual, 0)
		err = counter.Increment(ctx, 1)
		// if increment has error like shards has not been created, it will not have problem here, the error will be throw when commit transaction
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldNotBeNil)

	err = counters.DeleteSampleTotal(ctx)
	So(err, ShouldBeNil)
}

func testCounter(ctx context.Context, db SampleDB, counters *SampleCounters) {
	// clean counter
	err := counters.DeleteSampleTotal(ctx)
	So(err, ShouldBeNil)

	// create counter
	counter := counters.SampleTotal(ctx)
	So(counter, ShouldNotBeNil)
	So((counter.(*CounterFirestore)).errorID(), ShouldNotBeEmpty)
	err = counter.CreateShards(ctx)
	So(err, ShouldBeNil)

	//counter minimal shards is 10
	counter = counters.Counter("minShards", 0)
	So(counter, ShouldNotBeNil)
	firestoreCounter := counter.(*CounterFirestore)
	So(firestoreCounter.numShards, ShouldEqual, 10)
	err = counters.Delete(ctx, "minShards")
	So(err, ShouldBeNil)

	// delete exist counter
	err = counters.DeleteSampleTotal(ctx)
	So(err, ShouldBeNil)

	// new counter
	counter = counters.SampleTotal(ctx)
	So(counter, ShouldNotBeNil)
	err = counter.CreateShards(ctx)
	So(err, ShouldBeNil)

	count, err := counter.Count(ctx)
	So(count, ShouldEqual, 0)
	So(err, ShouldBeNil)

	err = counter.Increment(ctx, 2)
	So(err, ShouldBeNil)

	count, err = counter.Count(ctx)
	So(count, ShouldEqual, 2)
	So(err, ShouldBeNil)

	// exist counter
	counter2 := counters.SampleTotal(ctx)
	So(counter2, ShouldNotBeNil)

	count2, err := counter.Count(ctx)
	So(count2, ShouldEqual, 2)
	So(err, ShouldBeNil)

	err = counter.Increment(ctx, -2)
	So(err, ShouldBeNil)

	count2, err = counter.Count(ctx)
	So(count2, ShouldEqual, 0)
	So(err, ShouldBeNil)

	// get exist counter
	counter2 = counters.SampleTotal(ctx)
	So(counter2, ShouldNotBeNil)

	err = counter.Increment(ctx, 1)
	So(err, ShouldBeNil)

	count3, err := counter.Count(ctx)
	So(count3, ShouldEqual, 1)
	So(err, ShouldBeNil)

	//clean counter
	err = counters.DeleteSampleTotal(ctx)
	So(err, ShouldBeNil)

	//delete second time should be fine
	err = counters.DeleteSampleTotal(ctx)
	So(err, ShouldBeNil)

	err = counters.DeleteSampleTotal(ctx)
	So(err, ShouldBeNil)
}

func testCounterInCanceledCtx(ctx context.Context, db SampleDB, counters *SampleCounters) {

	counter := counters.SampleTotal(ctx)
	So(counter, ShouldNotBeNil)

	ctxCanceled := util.CanceledCtx()
	err := counter.CreateShards(ctxCanceled)
	So(err, ShouldNotBeNil)

	err = counter.Increment(ctxCanceled, 1)
	So(err, ShouldNotBeNil)

	count, err := counter.Count(ctxCanceled)
	So(err, ShouldNotBeNil)
	So(count, ShouldEqual, 0)

	err = counters.DeleteSampleTotal(ctx)
	So(err, ShouldBeNil)

	//test empty table name and empty counter name
	fCounter := counter.(*CounterFirestore)
	fCounter.counterName = ""
	count, err = fCounter.Count(ctx)
	So(err, ShouldNotBeNil)
	So(count, ShouldEqual, 0)
	err = fCounter.Increment(ctx, 1)
	So(err, ShouldNotBeNil)

	fCounter.tableName = ""
	count, err = fCounter.Count(ctx)
	So(err, ShouldNotBeNil)
	So(count, ShouldEqual, 0)
	err = fCounter.Increment(ctx, 1)
	So(err, ShouldNotBeNil)

}

func testCounterInTransaction(ctx context.Context, db SampleDB, counters *SampleCounters) {
	// clean counter
	err := counters.DeleteSampleTotal(ctx)
	So(err, ShouldBeNil)

	// do not read after write
	err = db.Transaction(ctx, func(ctx context.Context) error {
		counters := db.Counters()
		counter := counters.SampleTotal(ctx)
		So(counter, ShouldNotBeNil)
		err := counter.CreateShards(ctx)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		counters := db.Counters()
		counter := counters.SampleTotal(ctx)
		count, err := counter.Count(ctx)
		So(count, ShouldEqual, 0)
		So(err, ShouldBeNil)

		err = counter.Increment(ctx, 1)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)

	counter := counters.SampleTotal(ctx)
	count, err := counter.Count(ctx)
	So(count, ShouldEqual, 1)
	So(err, ShouldBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		counters := db.Counters()
		err = counters.DeleteSampleTotal(ctx)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)
	err = counters.DeleteSampleTotal(ctx)
	So(err, ShouldBeNil)
}

func TestConcurrentCounter(t *testing.T) {
	ctx := context.Background()
	db, _ := NewSampleGlobalDB(ctx)
	defer db.Close()
	counters := db.Counters()
	counters.DeleteSampleTotal(ctx)
	counter := counters.SampleTotal(ctx)
	counter.CreateShards(ctx)

	var concurrent = 3
	var wg sync.WaitGroup
	wg.Add(concurrent)
	counting := func() {
		counter := counters.SampleTotal(ctx)
		for i := 0; i < 5; i++ {
			err := counter.Increment(ctx, 1)
			if err != nil {
				t.Errorf("err should be nil, got %v", err)
				return
			}
			//			fmt.Printf("count:%v\n", i+1)
		}
		wg.Done()
	}

	//create 10 go routing to do counting
	for i := 0; i < concurrent; i++ {
		go counting()
	}
	wg.Wait()
	count, err := counter.Count(ctx)
	if err != nil {
		t.Errorf("err should be nil, got %v", err)
		return
	}
	if count != 15 {
		t.Errorf("count = %d; want 15", count)
	}
	//fmt.Printf("total count:%v\n", count)
}
