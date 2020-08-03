package usage

import (
	"context"
	"testing"
	"time"

	"github.com/piyuo/libsrv/data"

	. "github.com/smartystreets/goconvey/convey"
)

func NewSample(ctx context.Context) data.DB {
	conn, err := data.FirestoreGlobalConnection(ctx)
	if err != nil {
		panic(err)
	}
	db := &data.BaseDB{
		Conn: conn,
	}
	removeSample(ctx, db)
	return db
}

func removeSample(ctx context.Context, db data.DB) {
	table := &data.Table{
		Connection: db.Connection(),
		TableName:  "Usage",
	}
	table.Clear(ctx)
}

func TestUsage(t *testing.T) {
	Convey("Should count,add and remove usage", t, func() {
		ctx := context.Background()
		db := NewSample(ctx)
		defer removeSample(ctx, db)

		group := "test"
		key := "key1"
		usage := NewUsage(db)

		//check count is 0
		expired := time.Now().UTC().Add(time.Duration(-5) * time.Second)
		count, recent, err := usage.Count(ctx, group, key, expired)
		So(err, ShouldBeNil)
		So(count, ShouldEqual, 0)
		So(recent.IsZero(), ShouldBeTrue)

		//add usage
		err = usage.Add(ctx, group, key)
		So(err, ShouldBeNil)

		//check count is 1
		count, recent, err = usage.Count(ctx, group, key, expired)
		So(err, ShouldBeNil)
		So(count, ShouldEqual, 1)
		So(recent.IsZero(), ShouldBeFalse)
		//		dur := time.Now().UTC().Sub(recent)
		//		sec := dur.Seconds()
		So(time.Now().UTC().After(recent), ShouldBeTrue)

		//remove usage
		err = usage.Remove(ctx, group, key)
		So(err, ShouldBeNil)

		//check count is 0
		count, recent, err = usage.Count(ctx, group, key, expired)
		So(err, ShouldBeNil)
		So(count, ShouldEqual, 0)
		So(recent.IsZero(), ShouldBeTrue)
	})
}

func TestUsageDuration(t *testing.T) {
	Convey("Should only count usage in certain duration", t, func() {
		ctx := context.Background()
		db := NewSample(ctx)
		defer removeSample(ctx, db)

		group := "test"
		key := "key1"
		usage := NewUsage(db)

		//add usage that won't count
		err := usage.Add(ctx, group, key)
		So(err, ShouldBeNil)

		time.Sleep(time.Duration(2) * time.Second)

		err = usage.Add(ctx, group, key)
		So(err, ShouldBeNil)

		//check count is 1
		expired := time.Now().UTC().Add(time.Duration(-1) * time.Second)
		count, recent, err := usage.Count(ctx, group, key, expired)
		So(err, ShouldBeNil)
		So(count, ShouldEqual, 1)
		So(recent.IsZero(), ShouldBeFalse)
		So(time.Now().UTC().After(recent), ShouldBeTrue)
		dur := time.Now().UTC().Sub(recent)
		ms := dur.Milliseconds()
		So(ms < 1000, ShouldBeTrue)

		//remove usage
		err = usage.Remove(ctx, group, key)
		So(err, ShouldBeNil)
	})
}

func TestUsageMaintenance(t *testing.T) {
	Convey("Should maintenance usage", t, func() {
		ctx := context.Background()
		db := NewSample(ctx)
		defer removeSample(ctx, db)

		group := "test"
		key := "key1"
		usage := NewUsage(db)

		//add 2 usage
		err := usage.Add(ctx, group, key)
		So(err, ShouldBeNil)

		err = usage.Add(ctx, group, key)
		So(err, ShouldBeNil)

		time.Sleep(time.Duration(2) * time.Second)

		// test maintenance usage by remove past 1 seconds usage
		expired := time.Now().UTC().Add(time.Duration(-1) * time.Second)
		result, err := usage.Maintenance(ctx, expired)
		So(err, ShouldBeNil)
		So(result, ShouldBeTrue)

		count, _, err := usage.Count(ctx, group, key, expired)
		So(err, ShouldBeNil)
		So(count, ShouldEqual, 0)
	})
}
