package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSerial(t *testing.T) {
	Convey("should query table", t, func() {
		ctx := context.Background()
		dbG, dbR, samplesG, samplesR := firestoreBeginTest()
		defer dbG.Close()
		defer dbR.Close()
		serialG := dbG.Serial()
		serialR := dbR.Serial()

		serialTest(ctx, serialG)
		serialTest(ctx, serialR)

		firestoreEndTest(dbG, dbR, samplesG, samplesR)
	})
}

func serialTest(ctx context.Context, serial *SampleSerial) {
	err := serial.Delete(ctx, "sample-id")
	So(err, ShouldBeNil)

	num, err := serial.Number(ctx, "sample-id")
	So(err, ShouldBeNil)
	So(num, ShouldEqual, 1)

	num, err = serial.Number(ctx, "sample-id")
	So(err, ShouldBeNil)
	So(num, ShouldEqual, 2)

	num, err = serial.Number(ctx, "sample-id")
	So(err, ShouldBeNil)
	So(num, ShouldEqual, 3)

	id, err := serial.SampleID(ctx)
	So(err, ShouldBeNil)
	So(id, ShouldNotBeEmpty)

	err = serial.Delete(ctx, "sample-id")
	So(err, ShouldBeNil)
}

func emptyTableName(ctx context.Context, db SampleDB) {
	serial := db.Serial()
	serial.TableName = ""
	So(serial.TableName, ShouldBeEmpty)

	number, err := serial.Number(ctx, "sample-id")
	So(err, ShouldNotBeNil)
	So(number, ShouldEqual, 0)
	code, err := serial.Code(ctx, "sample-id")
	So(err, ShouldNotBeNil)
	So(code, ShouldBeEmpty)
}
