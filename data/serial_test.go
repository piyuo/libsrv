package data

import (
	"context"
	"testing"

	util "github.com/piyuo/libsrv/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSerial(t *testing.T) {
	Convey("should query table", t, func() {
		ctx := context.Background()
		dbG, dbR, samplesG, samplesR := firestoreBeginTest()
		defer dbG.Close()
		defer dbR.Close()

		serialTest(ctx, dbG)
		serialTest(ctx, dbR)
		serialEmptyTableNameTest(ctx, dbG)
		serialEmptyTableNameTest(ctx, dbR)
		serialInTransactionTest(ctx, dbG)
		serialInTransactionTest(ctx, dbR)
		serialContextCanceled(dbG)
		serialContextCanceled(dbR)
		firestoreEndTest(dbG, dbR, samplesG, samplesR)
	})
}

func serialTest(ctx context.Context, db SampleDB) {
	serial := db.Serial()
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

func serialInTransactionTest(ctx context.Context, db SampleDB) {
	serial := db.Serial()
	err := serial.Delete(ctx, "sample-id")

	err = db.Transaction(ctx, func(ctx context.Context) error {
		serial := db.Serial()
		num, err := serial.Number(ctx, "sample-id")
		So(err, ShouldBeNil)
		So(num, ShouldEqual, 1)
		return nil
	})
	So(err, ShouldBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		serial := db.Serial()
		num, err := serial.Number(ctx, "sample-id")
		So(err, ShouldBeNil)
		So(num, ShouldEqual, 2)
		// this line will case read after write error, only one serial can generated in transaction
		num, err = serial.Number(ctx, "sample-id")
		So(err, ShouldNotBeNil)
		So(num, ShouldEqual, 0)
		return nil
	})
	So(err, ShouldNotBeNil)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := serial.Delete(ctx, "sample-id")
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)
}

func serialEmptyTableNameTest(ctx context.Context, db SampleDB) {
	serial := db.Serial()
	serial.TableName = ""
	So(serial.TableName, ShouldBeEmpty)

	number, err := serial.Number(ctx, "sample-id")
	So(err, ShouldNotBeNil)
	So(number, ShouldEqual, 0)
	code, err := serial.Code(ctx, "sample-id")
	So(err, ShouldNotBeNil)
	So(code, ShouldBeEmpty)
	err = serial.Delete(ctx, "sample-id")
	So(err, ShouldNotBeNil)
}

func serialContextCanceled(db SampleDB) {
	ctx := util.CanceledCtx()

	serial := db.Serial()
	num, err := serial.Number(ctx, "sample-id")
	So(err, ShouldNotBeNil)
	So(num, ShouldEqual, 0)

	err = serial.Delete(ctx, "sample-id")
	So(err, ShouldNotBeNil)
}
