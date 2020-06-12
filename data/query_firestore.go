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

// NewQueryFirestore provide query for google firestore
func NewQueryFirestore(ctx context.Context, q firestore.Query, f func() Object) *QueryFirestore {
	return &QueryFirestore{
		AbstractQuery: AbstractQuery{ctx: ctx, newObject: f},
		query:         q}
}

//Where implement where on firestore
func (qf *QueryFirestore) Where(path, op string, value interface{}) Query {
	qf.query = qf.query.Where(path, op, value)
	return qf
}

//OrderBy implement orderby on firestore
func (qf *QueryFirestore) OrderBy(path string) Query {
	qf.query = qf.query.OrderBy(path, firestore.Asc)
	return qf
}

//OrderByDesc implement orderby desc on firestore
func (qf *QueryFirestore) OrderByDesc(path string) Query {
	qf.query = qf.query.OrderBy(path, firestore.Desc)
	return qf
}

//Limit implement limit on firestore
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

//Run query with default limit 100 object, use Limit() to override default limit
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
		obj := qf.newObject()
		err = doc.DataTo(obj)
		if err != nil {
			return err
		}
		obj.SetID(doc.Ref.ID)
		callback(obj)
	}
	return nil
}
