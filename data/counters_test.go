package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCounters(t *testing.T) {
	Convey("should check table name & counter name", t, func() {
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		cg, cr := createSampleCounters(dbG, dbR)
		defer removeSampleCounters(cg, cr)

		countersTest(t, cg)
		countersTest(t, cr)
	})

}

func countersTest(t *testing.T, counters *SampleCounters) {
	ctx := context.Background()

	counter := counters.Counter("sample-counter", 3)
	So(counter, ShouldNotBeNil)

	err := counter.CreateShards(ctx)
	So(err, ShouldBeNil)

	err = counter.Increment(ctx, 1)
	So(err, ShouldBeNil)

	count, err := counter.Count(ctx)
	So(err, ShouldBeNil)
	So(count, ShouldEqual, 1)
}
