package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSelectQuery(t *testing.T) {
	greet1 := &Greet{
		From:        "1",
		Description: "1",
	}
	greet2 := &Greet{
		From:        "2",
		Description: "2",
	}

	ctx := context.Background()
	db, _ := firestoreGlobalDB(ctx)
	defer db.Close()

	db.Put(ctx, greet1)
	db.Put(ctx, greet2)

	list, err := db.Select(ctx, GreetFactory).Where("From", "==", "1").Execute()
	err = db.Delete(ctx, list[0])
	Convey("delete select document", t, func() {
		So(err, ShouldBeNil)
	})

	Convey("test select count", t, func() {
		So(len(list), ShouldBeGreaterThanOrEqualTo, 1)
	})

	db.DeleteAll(ctx, greet1.ModelName(), 9)
}

func TestOrder(t *testing.T) {
	greet1 := &Greet{
		From:        "1",
		Description: "1",
	}
	greet2 := &Greet{
		From:        "2",
		Description: "2",
	}
	ctx := context.Background()
	db, _ := firestoreGlobalDB(ctx)
	defer db.Close()
	db.DeleteAll(ctx, greet1.ModelName(), 9)

	db.Put(ctx, greet1)
	db.Put(ctx, greet2)

	Convey("should Execute", t, func() {
		list, err := db.Select(ctx, GreetFactory).Execute()
		So(err, ShouldBeNil)
		So(len(list), ShouldEqual, 2)
	})

	Convey("OrderByDesc should return 2 first", t, func() {
		list, err := db.Select(ctx, GreetFactory).OrderByDesc("From").Execute()
		greet := list[0].(*Greet)
		So(err, ShouldBeNil)
		So(greet.From, ShouldEqual, "2")
	})

	Convey("OrderByDesc should return 1 first", t, func() {
		list, err := db.Select(ctx, GreetFactory).OrderBy("From").Execute()
		greet := list[0].(*Greet)
		So(err, ShouldBeNil)
		So(greet.From, ShouldEqual, "1")
	})

	Convey("Limit should return only 1 object", t, func() {
		list, err := db.Select(ctx, GreetFactory).Limit(1).Execute()
		So(err, ShouldBeNil)
		So(len(list), ShouldEqual, 1)
	})

	db.DeleteAll(ctx, GreetModelName, 9)
}

func TestStartEndAt(t *testing.T) {
	greet1 := &Greet{
		From: "a city",
	}
	greet2 := &Greet{
		From: "b city",
	}
	greet3 := &Greet{
		From: "c city",
	}
	ctx := context.Background()
	db, _ := firestoreGlobalDB(ctx)
	defer db.Close()
	db.DeleteAll(ctx, greet1.ModelName(), 9)
	db.Put(ctx, greet1)
	db.Put(ctx, greet2)
	db.Put(ctx, greet3)

	Convey("StartAt should start at", t, func() {
		list, err := db.Select(ctx, GreetFactory).OrderBy("From").StartAt("b city").Execute()
		So(err, ShouldBeNil)
		greet := list[0].(*Greet)
		So(greet.From, ShouldEqual, "b city")
		So(len(list), ShouldEqual, 2)
	})

	Convey("StartAt should start after", t, func() {
		list, err := db.Select(ctx, GreetFactory).OrderBy("From").StartAfter("b city").Execute()
		So(err, ShouldBeNil)
		greet := list[0].(*Greet)
		So(greet.From, ShouldEqual, "c city")
		So(len(list), ShouldEqual, 1)
	})

	Convey("EndAt should end at", t, func() {
		list, err := db.Select(ctx, GreetFactory).OrderBy("From").EndAt("b city").Execute()
		So(err, ShouldBeNil)
		greet := list[0].(*Greet)
		So(greet.From, ShouldEqual, "a city")
		So(len(list), ShouldEqual, 2)
	})

	Convey("EndAt should end before", t, func() {
		list, err := db.Select(ctx, GreetFactory).OrderBy("From").EndBefore("b city").Execute()
		So(err, ShouldBeNil)
		greet := list[0].(*Greet)
		So(greet.From, ShouldEqual, "a city")
		So(len(list), ShouldEqual, 1)
	})

	Convey("StartAt should start at and limit", t, func() {
		list, err := db.Select(ctx, GreetFactory).OrderBy("From").Limit(1).StartAt("b city").Execute()
		So(err, ShouldBeNil)
		greet := list[0].(*Greet)
		So(greet.From, ShouldEqual, "b city")
		So(len(list), ShouldEqual, 1)
	})

	db.DeleteAll(ctx, GreetModelName, 9)
}
