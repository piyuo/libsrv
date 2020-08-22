package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBatch(t *testing.T) {
	Convey("should batch work", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)
		samplesG, samplesR := createSampleTable(dbG, dbR)
		defer removeSampleTable(samplesG, samplesR)

		batchDeleteObjectTest(ctx, dbG, samplesG)
		batchDeleteTest(ctx, dbG, samplesG)
	})
}

func batchDeleteObjectTest(ctx context.Context, db SampleDB, table *Table) {
	sample1 := &Sample{
		Name:  "sample1",
		Value: 1,
	}
	sample2 := &Sample{
		Name:  "sample2",
		Value: 2,
	}

	count, err := table.Query().Count(ctx)
	So(count, ShouldEqual, 0)
	So(db.InBatch(), ShouldBeFalse)

	db.BatchBegin()
	So(db.InBatch(), ShouldBeTrue)
	table.Set(ctx, sample1) //batch mode do not return error
	table.Set(ctx, sample2)
	err = db.BatchCommit(ctx)
	So(err, ShouldBeNil)
	list, err := table.Query().Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 2)
	So(db.InBatch(), ShouldBeFalse)

	s1 := list[0].(*Sample)
	s2 := list[1].(*Sample)
	db.BatchBegin()
	table.Update(ctx, s1.ID, map[string]interface{}{
		"Value": 9,
	})
	table.Update(ctx, s2.ID, map[string]interface{}{
		"Value": 9,
	})
	err = db.BatchCommit(ctx)
	So(err, ShouldBeNil)

	list, err = table.Query().Where("Value", "==", 9).Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 2)
	g1 := list[0].(*Sample)
	g2 := list[1].(*Sample)
	So(g1.Value, ShouldEqual, 9)
	So(g2.Value, ShouldEqual, 9)

	db.BatchBegin()
	table.Increment(ctx, s1.ID, "Value", 1)
	table.Increment(ctx, s2.ID, "Value", 1)
	err = db.BatchCommit(ctx)
	So(err, ShouldBeNil)
	list, err = table.Query().Where("Value", "==", 10).Execute(ctx)
	So(err, ShouldBeNil)
	So(len(list), ShouldEqual, 2)
	i1 := list[0].(*Sample)
	i2 := list[1].(*Sample)
	So(i1.Value, ShouldEqual, 10)
	So(i2.Value, ShouldEqual, 10)

	db.BatchBegin()
	table.DeleteObject(ctx, sample1) //batch mode do not return error
	table.DeleteObject(ctx, sample2)
	err = db.BatchCommit(ctx)
	So(err, ShouldBeNil)
	count, err = table.Query().Count(ctx)
	So(err, ShouldBeNil)
	So(count, ShouldEqual, 0)
}

func batchDeleteTest(ctx context.Context, db SampleDB, table *Table) {
	sample1 := &Sample{
		Name:  "sample1",
		Value: 1,
	}
	sample2 := &Sample{
		Name:  "sample2",
		Value: 2,
	}

	db.BatchBegin()
	table.Set(ctx, sample1) //batch mode do not return error
	table.Set(ctx, sample2)
	err := db.BatchCommit(ctx)
	So(err, ShouldBeNil)

	idList, err := table.Query().GetIDs(ctx)
	So(err, ShouldBeNil)
	So(len(idList), ShouldEqual, 2)

	db.BatchBegin()
	table.Delete(ctx, idList[0]) //batch mode do not return error
	table.Delete(ctx, idList[1])
	err = db.BatchCommit(ctx)
	So(err, ShouldBeNil)
	count, err := table.Query().Count(ctx)
	So(err, ShouldBeNil)
	So(count, ShouldEqual, 0)
}
