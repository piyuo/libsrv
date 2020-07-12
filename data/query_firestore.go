package data

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// QueryFirestore implement google firestore
type QueryFirestore struct {
	Query
	query firestore.Query
	tx    *firestore.Transaction
}

// Where implement where on firestore
//
//	db.Select(ctx, GreetFactory).Where("From", "==", "1").Run(func(o Object) {
//		i++
//		err := db.Delete(ctx, o)
//	})
//
func (qf *QueryFirestore) Where(path, op string, value interface{}) QueryRef {
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
func (qf *QueryFirestore) OrderBy(path string) QueryRef {
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
func (qf *QueryFirestore) OrderByDesc(path string) QueryRef {
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
func (qf *QueryFirestore) Limit(n int) QueryRef {
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
func (qf *QueryFirestore) StartAt(docSnapshotOrFieldValues ...interface{}) QueryRef {
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
func (qf *QueryFirestore) StartAfter(docSnapshotOrFieldValues ...interface{}) QueryRef {
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
func (qf *QueryFirestore) EndAt(docSnapshotOrFieldValues ...interface{}) QueryRef {
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
func (qf *QueryFirestore) EndBefore(docSnapshotOrFieldValues ...interface{}) QueryRef {
	qf.query = qf.query.EndBefore(docSnapshotOrFieldValues...)
	return qf
}

// Execute query with default limit to 20 object, use Limit() to override default limit, return nil if anything wrong
//
//	list = []*Greet{}
//	ctx := context.Background()
//	db, _ := firestoreGlobalDB(ctx)
//	defer db.Close()
//	list, err := db.Select(ctx, GreetFactory).OrderBy("From").Limit(1).StartAt("b city").Execute(ctx)
//	greet := list[0].(*Greet)
//	So(greet.From, ShouldEqual, "b city")
//	So(len(list), ShouldEqual, 1)
//
func (qf *QueryFirestore) Execute(ctx context.Context) ([]ObjectRef, error) {
	if qf.limit == 0 {
		qf.Limit(limitQueryDefault)
	}
	var resultSet []ObjectRef

	var iter *firestore.DocumentIterator
	if qf.tx != nil {
		iter = qf.tx.Documents(qf.query)
	} else {
		iter = qf.query.Documents(ctx)
	}
	defer iter.Stop()

	for {
		snapshot, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		object := qf.factory()
		if object == nil {
			return nil, errors.New("failed to create object from factory")
		}

		err = snapshot.DataTo(object)
		if err != nil {
			return nil, err
		}
		object.SetRef(snapshot.Ref)
		object.SetID(snapshot.Ref.ID)
		object.SetCreateTime(snapshot.CreateTime)
		object.SetUpdateTime(snapshot.UpdateTime)
		object.SetReadTime(snapshot.ReadTime)
		resultSet = append(resultSet, object)
	}
	return resultSet, nil
}

// Count execute query and return max 10 count
//
//
func (qf *QueryFirestore) Count(ctx context.Context) (int, error) {
	if qf.limit == 0 {
		qf.Limit(limitQueryDefault)
	}
	var iter *firestore.DocumentIterator
	if qf.tx != nil {
		iter = qf.tx.Documents(qf.query)
	} else {
		iter = qf.query.Documents(ctx)
	}
	defer iter.Stop()

	count := 0
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, err
		}
		count++
	}
	return count, nil
}

// IsEmpty execute query and return true if no object exist in table
//
//
func (qf *QueryFirestore) IsEmpty(ctx context.Context) (bool, error) {
	qf.Limit(1)
	var iter *firestore.DocumentIterator
	if qf.tx != nil {
		iter = qf.tx.Documents(qf.query)
	} else {
		iter = qf.query.Documents(ctx)
	}
	defer iter.Stop()

	_, err := iter.Next()
	if err == iterator.Done {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return false, nil
}
