package data

import (
	"context"
	"testing"

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

		testCounter(ctx, counterG)
		testCounter(ctx, counterR)

		firestoreEndTest(dbG, dbR, samplesG, samplesR)
	})

}

func testCounter(ctx context.Context, counters *SampleCounters) {
	counter, err := counters.SampleTotal(ctx)
	So(counter, ShouldNotBeNil)
	So(err, ShouldBeNil)

	//counter minimal shards is 10
	counter, err = counters.Counter(ctx, "minShards", 0)
	So(counter, ShouldNotBeNil)
	So(err, ShouldBeNil)
	firestoreCounter := counter.(*CounterFirestore)
	So(firestoreCounter.N, ShouldEqual, 10)
	err = firestoreCounter.Delete(ctx)
	So(err, ShouldBeNil)

	// delete exist counter
	err = counter.Delete(ctx)
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
	err = counter2.Delete(ctx)
	So(err, ShouldBeNil)

	//delete second time should be fine
	err = counter.Delete(ctx)
	So(err, ShouldBeNil)

}
