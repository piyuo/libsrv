package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCounters(t *testing.T) {
	Convey("should check table name & counter name", t, func() {
		ctx := context.Background()
		dbG, dbR, samplesG, samplesR := firestoreBeginTest()
		countersG := dbG.Counters()
		countersG.TableName = ""
		err := countersG.DeleteSampleTotal(ctx)
		So(err, ShouldNotBeNil)

		countersR := dbR.Counters()
		err = countersR.Delete(ctx, "")
		So(err, ShouldNotBeNil)
		err = countersR.Delete(ctx, "")
		So(err, ShouldNotBeNil)

		defer dbG.Close()
		defer dbR.Close()

		firestoreEndTest(dbG, dbR, samplesG, samplesR)
	})

}
