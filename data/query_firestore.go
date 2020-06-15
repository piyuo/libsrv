package data

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// QueryFirestore implement google firestore
type QueryFirestore struct {
	AbstractQuery
	query firestore.Query
}

// NewQueryFirestore implement query on google firestore
//
//	obj := factory()
//	query := db.client.Collection(obj.ModelName()).Query
//	return NewQueryFirestore(ctx, query, factory)
//
func NewQueryFirestore(ctx context.Context, query firestore.Query, factory func() Object) *QueryFirestore {
	return &QueryFirestore{
		AbstractQuery: AbstractQuery{ctx: ctx, factory: factory},
		query:         query}
}

// Where implement where on firestore
//
//	db.Select(ctx, GreetFactory).Where("From", "==", "1").Run(func(o Object) {
//		i++
//		err := db.Delete(ctx, o)
//	})
//
func (qf *QueryFirestore) Where(path, op string, value interface{}) Query {
	qf.query = qf.query.Where(path, op, value)
	return qf
}

// OrderBy implement orderby on firestore
//
//	list = []*Greet{}
// 	db.Select(ctx, GreetFactory).OrderBy("From").Run(func(o Object) {
//		greet := o.(*Greet)
//		list = append(list, greet)
//	})
//
func (qf *QueryFirestore) OrderBy(path string) Query {
	qf.query = qf.query.OrderBy(path, firestore.Asc)
	return qf
}

// OrderByDesc implement orderby desc on firestore
//
//	list = []*Greet{}
// 	db.Select(ctx, GreetFactory).OrderByDesc("From").Run(func(o Object) {
//		greet := o.(*Greet)
//		list = append(list, greet)
//	})
//
func (qf *QueryFirestore) OrderByDesc(path string) Query {
	qf.query = qf.query.OrderBy(path, firestore.Desc)
	return qf
}

// Limit implement limit on firestore
//
//	list = []*Greet{}
//	db.Select(ctx, GreetFactory).Limit(1).Run(func(o Object) {
//		greet := o.(*Greet)
//		list = append(list, greet)
//	})
//
func (qf *QueryFirestore) Limit(n int) Query {
	qf.limit = n
	qf.query = qf.query.Limit(n)
	return qf
}

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
func (qf *QueryFirestore) StartAt(docSnapshotOrFieldValues ...interface{}) Query {
	qf.query = qf.query.StartAt(docSnapshotOrFieldValues...)
	return qf
}

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
func (qf *QueryFirestore) StartAfter(docSnapshotOrFieldValues ...interface{}) Query {
	qf.query = qf.query.StartAfter(docSnapshotOrFieldValues...)
	return qf
}

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
func (qf *QueryFirestore) EndAt(docSnapshotOrFieldValues ...interface{}) Query {
	qf.query = qf.query.EndAt(docSnapshotOrFieldValues...)
	return qf
}

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
func (qf *QueryFirestore) EndBefore(docSnapshotOrFieldValues ...interface{}) Query {
	qf.query = qf.query.EndBefore(docSnapshotOrFieldValues...)
	return qf
}

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
func (qf *QueryFirestore) Execute() ([]Object, error) {
	if qf.limit == 0 {
		qf.Limit(10)
	}
	var resultSet []Object
	iter := qf.query.Documents(qf.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		obj := qf.factory()
		err = doc.DataTo(obj)
		if err != nil {
			return nil, err
		}
		obj.SetID(doc.Ref.ID)
		resultSet = append(resultSet, obj)
	}
	return resultSet, nil
}
