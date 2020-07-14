package data

import (
	"context"
	"testing"

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

	})
}

func noErrorTest(ctx context.Context, table *Table) {
	So(table.Factory, ShouldNotBeNil)
	So(table.UUID(), ShouldNotBeEmpty)

	obj := table.NewObject()
	So(table.TableName, ShouldEqual, "sample")
	So(obj, ShouldNotBeNil)
	So((obj.(*Sample)).Name, ShouldBeEmpty)

	obj2 := table.Factory
	So(obj2, ShouldNotBeNil)

}
