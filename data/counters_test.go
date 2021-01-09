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

		countersTest(dbG, cg)
		countersTest(dbR, cr)
	})

}

func countersTest(db SampleDB, counters *SampleCounters) {
	ctx := context.Background()

	counter := counters.Counter("SampleCount", 3, DateHierarchyNone)
	So(counter, ShouldNotBeNil)

	err := counter.Clear(ctx)
	So(err, ShouldBeNil)
	defer counter.Clear(ctx)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(ctx)
		So(err, ShouldBeNil)
		return counter.IncrementWX(ctx, 1)
	})
	So(err, ShouldBeNil)

	count, err := counter.CountAll(ctx)
	So(err, ShouldBeNil)
	So(count, ShouldEqual, 1)
}
