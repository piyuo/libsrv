package data

import (
	"context"
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
		counterG := dbG.Counter()
		counterR := dbR.Counter()

		testCounter(ctx, dbG, counterG)
		testCounter(ctx, dbR, counterR)

		testCounterInCanceledCtx(ctx, dbG, counterG)
		testCounterInCanceledCtx(ctx, dbR, counterR)

		testCounterInTransaction(ctx, dbG, counterG)
		testCounterInTransaction(ctx, dbR, counterR)

		firestoreEndTest(dbG, dbR, samplesG, samplesR)
	})

}

func testCounter(ctx context.Context, db SampleDB, counters *SampleCounters) {
	// clean counter
	err := counters.DeleteSampleTotal(ctx)
	So(err, ShouldBeNil)

	// create counter
	counter, err := counters.SampleTotal(ctx)
	So(counter, ShouldNotBeNil)
	So(err, ShouldBeNil)
	So((counter.(*CounterFirestore)).errorID(), ShouldNotBeEmpty)

	//counter minimal shards is 10
	counter, err = counters.Counter(ctx, "minShards", 0)
	So(counter, ShouldNotBeNil)
	So(err, ShouldBeNil)
	firestoreCounter := counter.(*CounterFirestore)
	So(firestoreCounter.N, ShouldEqual, 10)
	err = counters.Delete(ctx, "minShards")
	So(err, ShouldBeNil)

	// delete exist counter
	err = counters.DeleteSampleTotal(ctx)
	So(err, ShouldBeNil)

	// new counter
	counter, err = counters.SampleTotal(ctx)
	So(counter, ShouldNotBeNil)
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
	counter2, err := counters.SampleTotal(ctx)
	So(counter2, ShouldNotBeNil)
	So(err, ShouldBeNil)

	count2, err := counter.Count(ctx)
	So(count2, ShouldEqual, 2)
	So(err, ShouldBeNil)

	err = counter.Increment(ctx, -2)
	So(err, ShouldBeNil)

	count2, err = counter.Count(ctx)
	So(count2, ShouldEqual, 0)
	So(err, ShouldBeNil)

	// get exist counter
	counter2, err = counters.SampleTotal(ctx)
	So(counter2, ShouldNotBeNil)
	So(err, ShouldBeNil)
	So(counter.GetCreateTime(), ShouldNotBeNil)
	So(counter.GetReadTime(), ShouldNotBeNil)
	So(counter.GetUpdateTime(), ShouldNotBeNil)

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
}

func testCounterInCanceledCtx(ctx context.Context, db SampleDB, counters *SampleCounters) {
	counter, err := counters.SampleTotal(ctx)
	So(counter, ShouldNotBeNil)
	So(err, ShouldBeNil)

	ctxCanceled := util.CanceledCtx()
	err = counter.Increment(ctxCanceled, 1)
	So(err, ShouldNotBeNil)

	count, err := counter.Count(ctxCanceled)
	So(err, ShouldNotBeNil)
	So(count, ShouldEqual, 0)

	err = counters.DeleteSampleTotal(ctx)
	So(err, ShouldBeNil)

}

func testCounterInTransaction(ctx context.Context, db SampleDB, counters *SampleCounters) {
	// clean counter
	err := counters.DeleteSampleTotal(ctx)
	So(err, ShouldBeNil)

	// do not read after write
	err = db.Transaction(ctx, func(ctx context.Context) error {
		counters := db.Counter()
		counter, err := counters.SampleTotal(ctx)
		So(err, ShouldBeNil)
		So(counter, ShouldNotBeNil)
		So(counter.GetCreateTime(), ShouldNotBeNil)
		So(counter.GetReadTime(), ShouldNotBeNil)
		So(counter.GetUpdateTime(), ShouldNotBeNil)
		return nil
	})
	So(err, ShouldBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		counters := db.Counter()
		counter, err := counters.SampleTotal(ctx)
		So(err, ShouldBeNil)
		count, err := counter.Count(ctx)
		So(count, ShouldEqual, 0)
		So(err, ShouldBeNil)

		err = counter.Increment(ctx, 1)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)

	counter, err := counters.SampleTotal(ctx)
	So(err, ShouldBeNil)
	count, err := counter.Count(ctx)
	So(count, ShouldEqual, 1)
	So(err, ShouldBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		counters := db.Counter()
		err = counters.DeleteSampleTotal(ctx)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)

}
