package data

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/session"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTable(t *testing.T) {
	Convey("should have no error", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		tableG, tableR := createSampleTable(dbG, dbR)
		defer removeSampleTable(tableG, tableR)

		noErrorTest(ctx, tableG)
		noErrorTest(ctx, tableR)

		searchTest(ctx, tableG)
		firstObjectTest(ctx, tableG)
	})
}

func noErrorTest(ctx context.Context, table *Table) {
	So(table.Factory, ShouldNotBeNil)
	So(table.UUID(), ShouldNotBeEmpty)

	obj := table.NewObject()
	So(table.TableName, ShouldEqual, "Sample")
	So(obj, ShouldNotBeNil)
	So((obj.(*Sample)).Name, ShouldBeEmpty)

	obj2 := table.Factory
	So(obj2, ShouldNotBeNil)
}

func firstObjectTest(ctx context.Context, table *Table) {
	sample1 := &Sample{
		Name:  "a",
		Value: 1,
	}
	table.Set(ctx, sample1)

	obj, err := table.GetFirstObject(ctx)
	So(err, ShouldBeNil)
	So(obj, ShouldNotBeNil)

	id, err := table.GetFirstID(ctx)
	So(err, ShouldBeNil)
	So(id, ShouldNotBeEmpty)

	err = table.Delete(ctx, id)
	So(err, ShouldBeNil)
}

func searchTest(ctx context.Context, table *Table) {

	sample1 := &Sample{
		Name:  "a",
		Value: 1,
	}
	sample2 := &Sample{
		Name:  "a",
		Value: 2,
	}
	table.Set(ctx, sample1)
	table.Set(ctx, sample2)

	list, err := table.SortList(ctx, "Name", "==", "a", "Value", DESC)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 2)
	obj1 := list[0].(*Sample)
	obj2 := list[1].(*Sample)
	So(obj1.Value, ShouldEqual, 2)
	So(obj2.Value, ShouldEqual, 1)

	list, err = table.SortList(ctx, "Name", "==", "a", "Value", ASC)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 2)
	obj1 = list[0].(*Sample)
	obj2 = list[1].(*Sample)
	So(obj1.Value, ShouldEqual, 1)
	So(obj2.Value, ShouldEqual, 2)
	table.Delete(ctx, obj1.ID)
	table.Delete(ctx, obj2.ID)
}

func TestChangedBy(t *testing.T) {
	Convey("should set by which user", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		tableG, tableR := createSampleTable(dbG, dbR)
		defer removeSampleTable(tableG, tableR)
		sample := &Sample{
			Name:  "a",
			Value: 1,
		}
		tableG.Set(ctx, sample)
		So(sample.GetBy(), ShouldBeEmpty)

		ctx = session.SetUserID(ctx, "user1")
		tableG.Set(ctx, sample)
		So(sample.GetBy(), ShouldEqual, "user1")

	})
}
