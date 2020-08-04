package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestQuery(t *testing.T) {
	Convey("should query table", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		tableG, tableR := createSampleTable(dbG, dbR)
		defer removeSampleTable(tableG, tableR)

		executeTopOneTest(ctx, tableG)
		listTest(ctx, tableG)
		queryTest(ctx, tableG)
	})
}

func queryTest(ctx context.Context, table *Table) {
	sample1 := &Sample{
		Name:  "sample1",
		Value: 1,
	}
	sample2 := &Sample{
		Name:  "sample2",
		Value: 2,
	}
	err := table.Set(ctx, sample1)
	So(err, ShouldBeNil)
	err = table.Set(ctx, sample2)
	So(err, ShouldBeNil)

	// get full object
	list, err := table.Query().Where("Name", "==", "sample1").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")

	// factory has no object return must error
	bakFactory := table.Factory
	table.Factory = func() Object {
		return nil
	}
	listX, err := table.Query().Where("Name", "==", "sample1").Execute(ctx)
	So(err, ShouldNotBeNil)
	So(listX, ShouldBeNil)
	table.Factory = bakFactory

	list, err = table.Query().Where("Name", "==", "sample2").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")

	list, err = table.Query().Where("Value", "==", 1).Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")

	list, err = table.Query().Where("Value", "==", 2).Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")

	//OrderBy,OrderByDesc
	list, err = table.Query().OrderBy("Name").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 2)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")

	list, err = table.Query().OrderByDesc("Name").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 2)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")

	//limit
	list, err = table.Query().OrderBy("Name").Limit(1).Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")

	list, err = table.Query().OrderByDesc("Name").Limit(1).Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")

	//startAt,startAfter,endAt,endBefore
	list, err = table.Query().OrderBy("Name").StartAt("sample2").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")

	list, err = table.Query().OrderBy("Name").StartAfter("sample1").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")

	list, err = table.Query().OrderBy("Name").EndAt("sample2").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 2)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")

	list, err = table.Query().OrderBy("Name").EndBefore("sample2").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")

	count, err := table.Query().Where("Name", "==", "sample1").Count(ctx)
	So(err, ShouldBeNil)
	So(count, ShouldEqual, 1)

	isEmpty, err := table.Query().Where("Name", "==", "sample1").IsEmpty(ctx)
	So(err, ShouldBeNil)
	So(isEmpty, ShouldBeFalse)

	isExist, err := table.Query().Where("Name", "==", "sample1").IsExist(ctx)
	So(err, ShouldBeNil)
	So(isExist, ShouldBeTrue)

	table.DeleteObject(ctx, sample1)
	table.DeleteObject(ctx, sample2)
}

func listTest(ctx context.Context, table *Table) {
	sample1 := &Sample{
		Name:  "sample",
		Value: 1,
	}
	sample2 := &Sample{
		Name:  "sample",
		Value: 2,
	}
	err := table.Set(ctx, sample1)
	So(err, ShouldBeNil)
	err = table.Set(ctx, sample2)
	So(err, ShouldBeNil)

	// get id only
	list, err := table.Query().Where("Name", "==", "sample").ExecuteID(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 2)
	So(list[0], ShouldNotBeEmpty)
	So(list[1], ShouldNotBeEmpty)
	So(list[0], ShouldNotEqual, list[1])

	table.DeleteObject(ctx, sample1)
	table.DeleteObject(ctx, sample2)

}

func executeTopOneTest(ctx context.Context, table *Table) {

	obj, err := table.Query().Where("Name", "==", "sample").ExecuteTopOne(ctx)
	So(err, ShouldBeNil)
	So(obj, ShouldBeNil)

	sample1 := &Sample{
		Name:  "sample",
		Value: 1,
	}
	sample2 := &Sample{
		Name:  "sample",
		Value: 2,
	}
	err = table.Set(ctx, sample1)
	So(err, ShouldBeNil)
	err = table.Set(ctx, sample2)
	So(err, ShouldBeNil)

	// get top one object only
	obj, err = table.Query().Where("Name", "==", "sample").ExecuteTopOne(ctx)
	So(err, ShouldBeNil)
	So(obj, ShouldNotBeNil)

	// set limit 2 still get 1 object
	obj, err = table.Query().Where("Name", "==", "sample").Limit(2).ExecuteTopOne(ctx)
	So(err, ShouldBeNil)
	So(obj, ShouldNotBeNil)

	table.DeleteObject(ctx, sample1)
	table.DeleteObject(ctx, sample2)

}
