package data

import (
	"context"
	"testing"

	util "github.com/piyuo/libsrv/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestShardsFirestore(t *testing.T) {
	Convey("check method", t, func() {
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

	Convey("check create and delete shards", t, func() {
		ctx := context.Background()
		dbG, dbR := createSampleDB()
		defer removeSampleDB(dbG, dbR)

		g := &ShardsFirestore{
			conn:      dbG.Connection.(*ConnectionFirestore),
			tableName: "sample-shards",
			id:        "sample-shard",
			numShards: 3,
		}

		r := &ShardsFirestore{
			conn:      dbR.Connection.(*ConnectionFirestore),
			tableName: "sample-shards",
			id:        "sample-shard",
			numShards: 3,
		}

		g.deleteShards(ctx)
		r.deleteShards(ctx)

		testShardsInTransaction(dbG, g)
		testShardsInTransaction(dbR, r)

		testShards(g)
		testShards(r)
	})

}

func testShards(shards *ShardsFirestore) {
	ctx := context.Background()
	docCount, shardsCount, err := shards.count(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 0)
	So(shardsCount, ShouldEqual, 0)

	err = shards.createShards(ctx)
	So(err, ShouldBeNil)

	docCount, shardsCount, err = shards.count(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 1)
	So(shardsCount, ShouldEqual, 3)

	// re-create shards with numShards=5
	shards.numShards = 5
	err = shards.createShards(ctx)
	So(err, ShouldBeNil)

	docCount, shardsCount, err = shards.count(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 1)
	So(shardsCount, ShouldEqual, 5)
	err = shards.deleteShards(ctx)
	So(err, ShouldBeNil)
	shards.numShards = 3

	docCount, shardsCount, err = shards.count(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 0)
	So(shardsCount, ShouldEqual, 0)

	//test ensureShardsDocument
	docRef, shardsRef := shards.getRef()
	shardRef := shardsRef.Doc("1")

	err = shards.ensureShardsDocument(ctx, docRef)
	So(err, ShouldBeNil)

	//try again
	err = shards.ensureShardsDocument(ctx, docRef)
	So(err, ShouldBeNil)

	docCount, shardsCount, err = shards.count(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 1)
	So(shardsCount, ShouldEqual, 0)

	err = shards.deleteShards(ctx)
	So(err, ShouldBeNil)

	//test ensureShard
	err = shards.ensureShard(ctx, docRef, shardRef)
	So(err, ShouldBeNil)

	//try again
	err = shards.ensureShard(ctx, docRef, shardRef)
	So(err, ShouldBeNil)

	docCount, shardsCount, err = shards.count(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 1)
	So(shardsCount, ShouldEqual, 1)

	err = shards.deleteShards(ctx)
	So(err, ShouldBeNil)
}

func testShardsInTransaction(db SampleDB, shards *ShardsFirestore) {
	ctx := context.Background()
	docCount, shardsCount, err := shards.count(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 0)
	So(shardsCount, ShouldEqual, 0)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		err = shards.createShards(ctx)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)

	docCount, shardsCount, err = shards.count(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 1)
	So(shardsCount, ShouldEqual, 3)

	// re-create shards with numShards=5
	shards.numShards = 5
	err = shards.createShards(ctx)
	So(err, ShouldBeNil)

	docCount, shardsCount, err = shards.count(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 1)
	So(shardsCount, ShouldEqual, 5)

	err = db.Transaction(ctx, func(ctx context.Context) error {
		err = shards.deleteShards(ctx)
		So(err, ShouldBeNil)
		return nil
	})
	So(err, ShouldBeNil)
	shards.numShards = 3

	docCount, shardsCount, err = shards.count(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 0)
	So(shardsCount, ShouldEqual, 0)

	//test ensureShardsDocument
	docRef, shardsRef := shards.getRef()
	shardRef := shardsRef.Doc("1")

	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := shards.ensureShardsDocument(ctx, docRef)
		So(err, ShouldBeNil)
		return nil
	})

	//try again
	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := shards.ensureShardsDocument(ctx, docRef)
		So(err, ShouldBeNil)
		return nil
	})

	docCount, shardsCount, err = shards.count(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 1)
	So(shardsCount, ShouldEqual, 0)

	err = shards.deleteShards(ctx)
	So(err, ShouldBeNil)

	//test ensureShard
	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := shards.ensureShard(ctx, docRef, shardRef)
		So(err, ShouldBeNil)
		return nil
	})

	//try again
	err = db.Transaction(ctx, func(ctx context.Context) error {
		err := shards.ensureShard(ctx, docRef, shardRef)
		So(err, ShouldBeNil)
		return nil
	})

	docCount, shardsCount, err = shards.count(ctx)
	So(err, ShouldBeNil)
	So(docCount, ShouldEqual, 1)
	So(shardsCount, ShouldEqual, 1)

	err = shards.deleteShards(ctx)
	So(err, ShouldBeNil)

}
