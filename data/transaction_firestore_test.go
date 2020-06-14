package data

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTransaction(t *testing.T) {
	greet1 := &Greet{
		From:        "1",
		Description: "1",
	}
	greet2 := &Greet{
		From:        "2",
		Description: "2",
	}
	ctx := context.Background()
	db, _ := NewGlobalDB(ctx)
	defer db.Close()

	db.DeleteAll(ctx, greet1.ModelName(), 9)
	Convey("transaction should not have error", t, func() {
		err := db.Transaction(ctx, func(ctx context.Context, tx Transaction) error {
			tx.Put(ctx, greet1)
			tx.Put(ctx, greet2)
			return nil
		})
		So(err, ShouldBeNil)
	})

	greet := &Greet{}
	greet.SetID(greet1.ID())
	Convey("transaction should get greet1", t, func() {
		err := db.Transaction(ctx, func(ctx context.Context, tx Transaction) error {
			tx.Get(ctx, greet)
			return nil
		})
		So(err, ShouldBeNil)
		So(greet.From, ShouldEqual, greet1.From)
	})

	Convey("transaction should delete object", t, func() {
		err := db.Transaction(ctx, func(ctx context.Context, tx Transaction) error {
			tx.Delete(ctx, greet)
			return nil
		})
		So(err, ShouldBeNil)
		exist, _ := db.ExistByID(ctx, GreetModelName, greet1.ID())
		So(exist, ShouldBeFalse)
	})

	Convey("transaction should list object", t, func() {
		list, _ := db.ListAll(ctx, GreetFactory, 100)
		So(len(list), ShouldEqual, 1)
	})

	db.DeleteAll(ctx, greet1.ModelName(), 9)
	Convey("transaction should have error", t, func() {
		err := db.Transaction(ctx, func(ctx context.Context, tx Transaction) error {
			tx.Put(ctx, greet1)
			return errors.New("some thing wrong")
		})
		So(err, ShouldNotBeNil)
	})
	Convey("transaction fail should rollback", t, func() {
		list, _ := db.ListAll(ctx, GreetFactory, 100)
		So(len(list), ShouldEqual, 0)
	})

	db.DeleteAll(ctx, greet1.ModelName(), 9)
}

func TestTransactionID(t *testing.T) {
	Convey("transaction should create unique id", t, func() {
		ctx := context.Background()
		db, _ := NewGlobalDB(ctx)
		defer db.Close()
		db.DeleteByID(ctx, "shortID", "myID")
		err := db.Transaction(ctx, func(ctx context.Context, tx Transaction) error {
			id, err := tx.ShortID(ctx, "myID")
			So(err, ShouldBeNil)
			So(id.Number(), ShouldEqual, 1)
			return nil
		})
		err = db.Transaction(ctx, func(ctx context.Context, tx Transaction) error {
			id, err := tx.ShortID(ctx, "myID")
			So(err, ShouldBeNil)
			So(id.Number(), ShouldEqual, 2)
			return nil
		})
		db.DeleteByID(ctx, "shortID", "myID")
		So(err, ShouldBeNil)
	})
}
