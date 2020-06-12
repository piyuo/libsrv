package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSelectQuery(t *testing.T) {
	greet1 := Greet{
		From:        "1",
		Description: "1",
	}
	greet2 := Greet{
		From:        "2",
		Description: "2",
	}

	ctx := context.Background()
	db, _ := NewGlobalDB(ctx)
	defer db.Close()

	db.Put(ctx, &greet1)
	db.Put(ctx, &greet2)

	//test select
	var i int
	db.Select(ctx, GreetFactory).Where("From", "==", "1").Run(func(o Object) {
		i++
		err := db.Delete(ctx, o)
		Convey("delete select document ", t, func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("test select count '", t, func() {
		So(i, ShouldBeGreaterThanOrEqualTo, 1)
	})

	db.DeleteAll(ctx, greet1.ModelName(), 9)
}

func TestOrder(t *testing.T) {
	greet1 := Greet{
		From:        "1",
		Description: "1",
	}
	greet2 := Greet{
		From:        "2",
		Description: "2",
	}
	ctx := context.Background()
	db, _ := NewGlobalDB(ctx)
	defer db.Close()

	db.Put(ctx, &greet1)
	db.Put(ctx, &greet2)

	list := []*Greet{}
	db.Select(ctx, GreetFactory).OrderByDesc("From").Run(func(o Object) {
		greet := o.(*Greet)
		list = append(list, greet)
	})
	Convey("OrderByDesc should return 2 first ", t, func() {
		So(list[0].From, ShouldEqual, "2")
	})

	list = []*Greet{}
	db.Select(ctx, GreetFactory).OrderBy("From").Run(func(o Object) {
		greet := o.(*Greet)
		list = append(list, greet)
	})
	Convey("OrderByDesc should return 1 first ", t, func() {
		So(list[0].From, ShouldEqual, "1")
	})

	list = []*Greet{}
	db.Select(ctx, GreetFactory).Limit(1).Run(func(o Object) {
		greet := o.(*Greet)
		list = append(list, greet)
	})
	Convey("Limit should return only 1 object ", t, func() {
		So(len(list), ShouldEqual, 1)
	})

	db.DeleteAll(ctx, greet1.ModelName(), 9)
}

/*
func TestOffset(t *testing.T) {
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
	db.Put(&greet1)
	db.Put(&greet2)

	list := []*Greet{}
	db.Select(GreetFactory()).OrderBy("From").Offset(1).Run(func(o Object) {
		greet := o.(*Greet)
		list = append(list, greet)
	})
	Convey("StartAt should return 1 object ", t, func() {
		So(len(list), ShouldEqual, 1)
	})
	Convey("return object is 2 ", t, func() {
		So(list[0].From, ShouldEqual, "2")
	})
	db.DeleteAll(greet1.Class(), 9)
}
*/
