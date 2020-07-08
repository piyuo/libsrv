package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTable(t *testing.T) {
	Convey("should have no error", t, func() {
		ctx := context.Background()
		dbG, dbR, samplesG, samplesR := firestoreBeginTest()
		defer dbG.Close()
		defer dbR.Close()

		noErrorTest(ctx, samplesG)
		noErrorTest(ctx, samplesR)

		firestoreEndTest(dbG, dbR, samplesG, samplesR)
	})
}

func noErrorTest(ctx context.Context, table *Table) {
	So(table.Factory, ShouldNotBeNil)

	obj := table.NewObject()
	So(table.TableName(), ShouldEqual, "sample")
	So(obj, ShouldNotBeNil)
	So((obj.(*Sample)).Name, ShouldBeEmpty)

	obj2 := table.Factory()
	So(obj2, ShouldNotBeNil)

	id := table.ID()
	So(id, ShouldNotBeEmpty)

}
