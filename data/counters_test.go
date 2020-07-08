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
		counterG := dbG.Counter()
		counterG.TableName = ""
		err := counterG.DeleteSampleTotal(ctx)
		So(err, ShouldNotBeNil)
		_, err = counterG.SampleTotal(ctx)
		So(err, ShouldNotBeNil)

		counterR := dbR.Counter()
		err = counterR.Delete(ctx, "")
		So(err, ShouldNotBeNil)
		err = counterR.Delete(ctx, "")
		So(err, ShouldNotBeNil)

		defer dbG.Close()
		defer dbR.Close()

		firestoreEndTest(dbG, dbR, samplesG, samplesR)
	})

}
