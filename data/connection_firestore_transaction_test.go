package data

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTransaction(t *testing.T) {
	Convey("should transaction work", t, func() {
		ctx := context.Background()
		dbG, dbR, samplesG, samplesR := firestoreBeginTest()
		defer dbG.Close()
		defer dbR.Close()

		transactionTest(ctx, dbG, samplesG)
		transactionTest(ctx, dbR, samplesR)

		firestoreEndTest(dbG, dbR, samplesG, samplesR)
	})
}

func transactionTest(ctx context.Context, db DB, table Table) {

	sample1 := &Sample{
		Name:  "sample1",
		Value: 1,
	}
	sample2 := &Sample{
		Name:  "sample2",
		Value: 2,
	}

	//success transaction
	err := db.Transaction(ctx, func(ctx context.Context) error {
		err := table.Set(ctx, sample1)
		So(err, ShouldBeNil)
		err = table.Set(ctx, sample2)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)

	list, err := table.Query(ctx).OrderBy("Name").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 2)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")
	So((list[1].(*Sample)).Name, ShouldEqual, "sample2")
	isEmpty, err := table.IsEmpty(ctx)
	So(isEmpty, ShouldBeFalse)
	err = table.Clear(ctx)
	So(err, ShouldBeNil)

	//fail transaction
	err = db.Transaction(ctx, func(ctx context.Context) error {
		err = table.Set(ctx, sample1)
		So(err, ShouldBeNil)
		return errors.New("something wrong")
	})
	So(err, ShouldNotBeNil)

	isEmpty, err = table.IsEmpty(ctx)
	So(isEmpty, ShouldBeTrue)

	// success delete
	err = db.Transaction(ctx, func(ctx context.Context) error {
		err = table.Set(ctx, sample1)
		So(err, ShouldBeNil)
		err = table.DeleteObject(ctx, sample1)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)

	isEmpty, err = table.IsEmpty(ctx)
	So(isEmpty, ShouldBeTrue)
	err = table.Clear(ctx)
	So(err, ShouldBeNil)

	// failed delete
	err = db.Transaction(ctx, func(ctx context.Context) error {
		err = table.Set(ctx, sample1)
		So(err, ShouldBeNil)
		err = table.DeleteObject(ctx, sample1)
		So(err, ShouldBeNil)
		return errors.New("something wrong")
	})
	So(err, ShouldNotBeNil)

	isEmpty, err = table.IsEmpty(ctx)
	So(isEmpty, ShouldBeTrue)
	err = table.Clear(ctx)
	So(err, ShouldBeNil)

}
