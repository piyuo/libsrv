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

//Offset implement start at on firestore, often use by paginate data
//in firestore will bill extra mony on offset
//func (q *QueryFirestore) Offset(n int) IQuery {
//	q.query = q.query.Offset(n)
//	return q
//}

// Run query with default limit 100 object, use Limit() to override default limit
//
//	list = []*Greet{}
//	db.Select(ctx, GreetFactory).Run(func(o Object) {
//		greet := o.(*Greet)
//		list = append(list, greet)
//	})
//
func (qf *QueryFirestore) Run(callback func(o Object)) error {
	if qf.limit == 0 {
		qf.Limit(100)
	}

	iter := qf.query.Documents(qf.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		obj := qf.factory()
		err = doc.DataTo(obj)
		if err != nil {
			return err
		}
		obj.SetID(doc.Ref.ID)
		callback(obj)
	}
	return nil
}
