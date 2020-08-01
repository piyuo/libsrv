package data

import (
	"context"
	"testing"

	util "github.com/piyuo/libsrv/util"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTransaction(t *testing.T) {
	Convey("should transaction work", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		samplesG, samplesR := createSampleTable(dbG, dbR)
		defer removeSampleTable(samplesG, samplesR)

		transactionTest(ctx, dbG, samplesG)
		transactionTest(ctx, dbR, samplesR)
		methodTest(ctx, dbG, samplesG, true)
		methodTest(ctx, dbR, samplesR, false)
	})
}

func transactionTest(ctx context.Context, db SampleDB, table *Table) {
	sample1 := &Sample{
		Name:  "sample1",
		Value: 1,
	}
	sample2 := &Sample{
		Name:  "sample2",
		Value: 2,
	}

	So(db.IsInTransaction(), ShouldBeFalse)
	//success transaction
	err := db.Transaction(ctx, func(ctx context.Context) error {
		So(db.IsInTransaction(), ShouldBeTrue)
		err := table.Set(ctx, sample1)
		So(err, ShouldBeNil)
		err = table.Set(ctx, sample2)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)

	list, err := table.Query().OrderBy("Name").Execute(ctx)
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

func methodTest(ctx context.Context, db SampleDB, table *Table, isGlobal bool) {

	sample1 := &Sample{
		Name:  "sample1",
		Value: 1,
	}
	sample2 := &Sample{
		Name:  "sample2",
		Value: 2,
	}

	// get & deleteObject
	err := table.Set(ctx, sample1)
	So(err, ShouldBeNil)
	err = db.Transaction(ctx, func(ctx context.Context) error {
		sample, err := table.Get(ctx, sample1.ID)
		So(err, ShouldBeNil)
		err = table.DeleteObject(ctx, sample)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)
	isEmpty, err := table.IsEmpty(ctx)
	So(isEmpty, ShouldBeTrue)

	// exist & list & delete
	err = table.Set(ctx, sample1)
	So(err, ShouldBeNil)
	err = db.Transaction(ctx, func(ctx context.Context) error {
		exist, err := table.Exist(ctx, sample1.ID)
		So(err, ShouldBeNil)
		So(exist, ShouldBeTrue)
		objects, err := table.All(ctx)
		So(err, ShouldBeNil)
		So(len(objects), ShouldEqual, 1)
		err = table.Delete(ctx, sample1.ID)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)
	isEmpty, err = table.IsEmpty(ctx)
	So(isEmpty, ShouldBeTrue)

	// select & update & Increment
	err = table.Set(ctx, sample1)
	So(err, ShouldBeNil)
	err = db.Transaction(ctx, func(ctx context.Context) error {
		name, err := table.Select(ctx, sample1.ID, "Name")
		So(err, ShouldBeNil)
		So(name.(string), ShouldEqual, "sample1")
		err = table.Update(ctx, sample1.ID, map[string]interface{}{
			"Name": "sample",
		})
		So(err, ShouldBeNil)
		err = table.Increment(ctx, sample1.ID, "Value", 1)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)
	name, err := table.Select(ctx, sample1.ID, "Name")
	So(err, ShouldBeNil)
	So(name.(string), ShouldEqual, "sample")
	value, err := table.Select(ctx, sample1.ID, "Value")
	So(err, ShouldBeNil)
	intValue, err := util.ToInt(value)
	So(err, ShouldBeNil)
	So(intValue, ShouldEqual, 2)
	table.DeleteObject(ctx, sample1)

	// query & clear
	err = table.Set(ctx, sample1)
	So(err, ShouldBeNil)
	err = table.Set(ctx, sample2)
	So(err, ShouldBeNil)
	err = db.Transaction(ctx, func(ctx context.Context) error {
		obj, err := table.Find(ctx, "Name", "==", "sample1")
		So(err, ShouldBeNil)
		So((obj.(*Sample)).Name, ShouldEqual, "sample1")

		list, err := table.Query().OrderBy("Name").Execute(ctx)
		So(err, ShouldBeNil)
		So(len(list), ShouldEqual, 2)
		So(list[0].(*Sample).Name, ShouldEqual, sample1.Name)
		So(list[1].(*Sample).Name, ShouldEqual, sample2.Name)

		err = table.Clear(ctx)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)
	obj, err := table.Find(ctx, "Value", "==", 2)
	So(err, ShouldBeNil)
	So(obj, ShouldBeNil)
	isEmpty, err = table.IsEmpty(ctx)
	So(isEmpty, ShouldBeTrue)

	// search & count & is empty
	err = table.Set(ctx, sample1)
	So(err, ShouldBeNil)
	err = db.Transaction(ctx, func(ctx context.Context) error {

		objects, err := table.List(ctx, "Name", "==", "sample1")
		So(err, ShouldBeNil)
		So(len(objects), ShouldEqual, 1)

		count, err := table.Count(ctx)
		So(err, ShouldBeNil)
		So(count, ShouldEqual, 1)

		empty, err := table.IsEmpty(ctx)
		So(err, ShouldBeNil)
		So(empty, ShouldEqual, false)

		err = table.DeleteObject(ctx, sample1)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)
	isEmpty, err = table.IsEmpty(ctx)
	So(isEmpty, ShouldBeTrue)

	//create & delete namespace
	err = db.Transaction(ctx, func(ctx context.Context) error {
		err = db.CreateNamespace(ctx)
		if isGlobal {
			So(err, ShouldNotBeNil)
		} else {
			So(err, ShouldBeNil)
		}
		err = db.DeleteNamespace(ctx)
		if isGlobal {
			So(err, ShouldNotBeNil)
		} else {
			So(err, ShouldBeNil)
		}
		return nil
	})
	So(err, ShouldBeNil)

}
