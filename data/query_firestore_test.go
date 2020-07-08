package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestQuery(t *testing.T) {
	Convey("should query table", t, func() {
		ctx := context.Background()
		dbG, dbR, samplesG, samplesR := firestoreBeginTest()
		defer dbG.Close()
		defer dbR.Close()

		queryTest(ctx, samplesG)
		queryTest(ctx, samplesR)

		firestoreEndTest(dbG, dbR, samplesG, samplesR)
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

	list, err := table.Query(ctx).Where("Name", "==", "sample1").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")

	list, err = table.Query(ctx).Where("Name", "==", "sample2").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")

	list, err = table.Query(ctx).Where("Value", "==", 1).Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")

	list, err = table.Query(ctx).Where("Value", "==", 2).Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")

	//OrderBy,OrderByDesc
	list, err = table.Query(ctx).OrderBy("Name").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 2)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")

	list, err = table.Query(ctx).OrderByDesc("Name").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 2)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")

	//limit
	list, err = table.Query(ctx).OrderBy("Name").Limit(1).Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")

	list, err = table.Query(ctx).OrderByDesc("Name").Limit(1).Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")

	//startAt,startAfter,endAt,endBefore
	list, err = table.Query(ctx).OrderBy("Name").StartAt("sample2").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")

	list, err = table.Query(ctx).OrderBy("Name").StartAfter("sample1").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")

	list, err = table.Query(ctx).OrderBy("Name").EndAt("sample2").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 2)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")

	list, err = table.Query(ctx).OrderBy("Name").EndBefore("sample2").Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 1)
	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")

}
