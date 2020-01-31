package data

import (
	"context"


	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// QueryFirestore implement google firestore
type QueryFirestore struct {
	Query
	query     firestore.Query
	ctx       context.Context
	newObject func() Object
	limit     int
}

// NewQueryFirestore provide query for google firestore
func NewQueryFirestore(ctx context.Context, q firestore.Query, f func() Object) *QueryFirestore {
	return &QueryFirestore{ctx: ctx, query: q, newObject: f}
}

//Where implement where on firestore
func (q *QueryFirestore) Where(path, op string, value interface{}) Query {
	q.query = q.query.Where(path, op, value)
	return q
}

//OrderBy implement orderby on firestore
func (q *QueryFirestore) OrderBy(path string) Query {
	q.query = q.query.OrderBy(path, firestore.Asc)
	return q
}

//OrderByDesc implement orderby desc on firestore
func (q *QueryFirestore) OrderByDesc(path string) Query {
	q.query = q.query.OrderBy(path, firestore.Desc)
	return q
}

//Limit implement limit on firestore
func (q *QueryFirestore) Limit(n int) Query {
	q.limit = n
	q.query = q.query.Limit(n)
	return q
}

//Offset implement start at on firestore, often use by paginate data
//in firestore will bill extra mony on offset
//func (q *QueryFirestore) Offset(n int) IQuery {
//	q.query = q.query.Offset(n)
//	return q
//}

//Run query with default limit 100 object, use Limit() to override default limit
func (q *QueryFirestore) Run(callback func(o Object)) error {

	if q.limit == 0 {
		q.Limit(100)
	}

	iter := q.query.Documents(q.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		obj := q.newObject()
		err = doc.DataTo(obj)
		if err != nil {
			return err
		}
		obj.SetID(doc.Ref.ID)
		callback(obj)
	}
	return nil
}
