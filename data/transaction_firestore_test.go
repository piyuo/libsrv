package data

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTransaction(t *testing.T) {
	greet1 := Greet{
		From:        "1",
		Description: "1",
	}
	greet2 := Greet{
		From:        "2",
		Description: "2",
	}
	ctx := context.Background()
	db, _ := firestoreNewDB(ctx)
	defer db.Close()

	db.DeleteAll(ctx, greet1.ModelName(), 9)
	err := db.RunTransaction(ctx, func(ctx context.Context, tx Transaction) error {
		tx.Put(ctx, &greet1)
		tx.Put(ctx, &greet2)
		return nil
	})
	Convey("transaction should not have error '", t, func() {
		So(err, ShouldBeNil)
	})
	list, _ := db.ListAll(ctx, GreetFactory, 100)
	Convey("transaction fail should rollback '", t, func() {
		So(len(list), ShouldEqual, 2)
	})

	db.DeleteAll(ctx, greet1.ModelName(), 9)
	err = db.RunTransaction(ctx, func(ctx context.Context, tx Transaction) error {
		tx.Put(ctx, &greet1)
		return errors.New("some thing wrong")
	})
	Convey("transaction should have error '", t, func() {
		So(err, ShouldNotBeNil)
	})
	list, _ = db.ListAll(ctx, GreetFactory, 100)
	Convey("transaction fail should rollback '", t, func() {
		So(len(list), ShouldEqual, 0)
	})

	db.DeleteAll(ctx, greet1.ModelName(), 9)
}
