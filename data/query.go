package data

import "context"

// Query is query interface
type Query interface {
	// Where set where filter
	//
	//	db.Select(ctx, GreetFactory).Where("From", "==", "1").Run(func(o Object) {
	//		i++
	//		err := db.Delete(ctx, o)
	//	})
	//
	Where(path, op string, value interface{}) Query

	// OrderBy set query order by
	//
	//	list = []*Greet{}
	// 	db.Select(ctx, GreetFactory).OrderBy("From").Run(func(o Object) {
	//		greet := o.(*Greet)
	//		list = append(list, greet)
	//	})
	//
	OrderBy(path string) Query

	// OrderByDesc set query order by desc
	//
	//	list = []*Greet{}
	// 	db.Select(ctx, GreetFactory).OrderByDesc("From").Run(func(o Object) {
	//		greet := o.(*Greet)
	//		list = append(list, greet)
	//	})
	//
	OrderByDesc(path string) Query

	// Limit set query limit
	//
	//	list = []*Greet{}
	//	db.Select(ctx, GreetFactory).Limit(1).Run(func(o Object) {
	//		greet := o.(*Greet)
	//		list = append(list, greet)
	//	})
	//
	Limit(n int) Query

	// Execute query with default limit to 10 object, use Limit() to override default limit, return nil if anything wrong
	//
	//	list = []*Greet{}
	//	ctx := context.Background()
	//	db, _ := firestoreGlobalDB(ctx)
	//	defer db.Close()
	//	list, err := db.Select(ctx, GreetFactory).OrderBy("From").Limit(1).StartAt("b city").Execute()
	//	greet := list[0].(*Greet)
	//	So(greet.From, ShouldEqual, "b city")
	//	So(len(list), ShouldEqual, 1)
	//
	Execute() ([]Object, error)

	// StartAt implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
	//
	//	greet1 := &Greet{
	//		From: "a city",
	//	}
	//	greet2 := &Greet{
	//		From: "b city",
	//	}
	//	list, err := db.Select(ctx, GreetFactory).OrderBy("From").StartAt("b city").Execute()
	//	So(err, ShouldBeNil)
	//	greet := list[0].(*Greet)
	//	So(greet.From, ShouldEqual, "b city")
	//	So(len(list), ShouldEqual, 2)
	//
	StartAt(docSnapshotOrFieldValues ...interface{}) Query

	// StartAfter implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
	//
	//	greet1 := &Greet{
	//		From: "a city",
	//	}
	//	greet2 := &Greet{
	//		From: "b city",
	//	}
	//	list, err := db.Select(ctx, GreetFactory).OrderBy("From").StartAfter("b city").Execute()
	//	So(err, ShouldBeNil)
	//	greet := list[0].(*Greet)
	//	So(greet.From, ShouldEqual, "c city")
	//	So(len(list), ShouldEqual, 1)
	//
	StartAfter(docSnapshotOrFieldValues ...interface{}) Query

	// EndAt implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
	//
	//	greet1 := &Greet{
	//		From: "a city",
	//	}
	//	greet2 := &Greet{
	//		From: "b city",
	//	}
	//	list, err := db.Select(ctx, GreetFactory).OrderBy("From").EndAt("b city").Execute()
	//	So(err, ShouldBeNil)
	//	greet := list[0].(*Greet)
	//	So(greet.From, ShouldEqual, "a city")
	//	So(len(list), ShouldEqual, 2)
	//
	EndAt(docSnapshotOrFieldValues ...interface{}) Query

	// EndBefore implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
	//
	//	greet1 := &Greet{
	//		From: "a city",
	//	}
	//	greet2 := &Greet{
	//		From: "b city",
	//	}
	//	list, err := db.Select(ctx, GreetFactory).OrderBy("From").EndBefore("b city").Execute()
	//	So(err, ShouldBeNil)
	//	greet := list[0].(*Greet)
	//	So(greet.From, ShouldEqual, "a city")
	//	So(len(list), ShouldEqual, 1)
	//
	EndBefore(docSnapshotOrFieldValues ...interface{}) Query
}

// AbstractQuery is query object need to implement
type AbstractQuery struct {
	Query

	// ctx is context
	//
	ctx context.Context

	// factor use to create object
	//
	factory func() Object

	// limit remember if query set limit, if not we will give default limit (10) to avoid return too may document
	//
	limit int
}
