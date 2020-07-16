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

		countersTest(dbG, cg)
		countersTest(dbR, cr)
	})

}

func countersTest(db SampleDB, counters *SampleCounters) {
	ctx := context.Background()

	counter := counters.Counter("sample-counter", 3)
	So(counter, ShouldNotBeNil)

	err := db.Transaction(ctx, func(ctx context.Context) error {
		err := counter.IncrementRX(1)
		So(err, ShouldBeNil)
		return counter.IncrementWX()
	})
	So(err, ShouldBeNil)

	count, err := counter.Count(ctx)
	So(err, ShouldBeNil)
	So(count, ShouldEqual, 1)
}
