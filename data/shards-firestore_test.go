package data

import (
	"context"
	"testing"

	util "github.com/piyuo/libsrv/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShardsFirestore(t *testing.T) {
	Convey("Should work normally", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)

		shards := ShardsFirestore{
			conn:      dbG.Connection.(*ConnectionFirestore),
			tableName: "tablename",
			id:        "id",
			numShards: 0,
		}
		id := shards.errorID()
		So(id, ShouldEqual, "tablename{root}-id")

		docRef, shardsRef := shards.getRef()
		So(docRef, ShouldNotBeNil)
		So(shardsRef, ShouldNotBeNil)

		//check canceled ctx
		ctxCanceled := util.CanceledCtx()
		err := shards.assert(ctxCanceled)
		So(err, ShouldNotBeNil)

		//check empty id
		shards.id = ""
		err = shards.assert(ctx)
		So(err, ShouldNotBeNil)

		//check empty table name
		shards.id = "id"
		shards.tableName = ""
		err = shards.assert(ctx)
		So(err, ShouldNotBeNil)

	})
}
