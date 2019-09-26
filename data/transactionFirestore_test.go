package data

import (
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
	db, _ := ProviderInstance().NewDB()
	defer db.Close()

	db.DeleteAll(greet1.Class(), 9)
	err := db.RunTransaction(func(tx ITransaction) error {
		tx.Put(&greet1)
		tx.Put(&greet2)
		return nil
	})
	Convey("transaction should not have error '", t, func() {
		So(err, ShouldBeNil)
	})
	list, _ := db.ListAll(GreetFactory, 100)
	Convey("transaction fail should rollback '", t, func() {
		So(len(list), ShouldEqual, 2)
	})

	db.DeleteAll(greet1.Class(), 9)
	err = db.RunTransaction(func(tx ITransaction) error {
		tx.Put(&greet1)
		return errors.New("some thing wrong")
	})
	Convey("transaction should have error '", t, func() {
		So(err, ShouldNotBeNil)
	})
	list, _ = db.ListAll(GreetFactory, 100)
	Convey("transaction fail should rollback '", t, func() {
		So(len(list), ShouldEqual, 0)
	})

	db.DeleteAll(greet1.Class(), 9)
}
