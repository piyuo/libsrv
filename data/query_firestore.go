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
func (class *QueryFirestore) Where(path, op string, value interface{}) Query {
	class.query = class.query.Where(path, op, value)
	return class
}

//OrderBy implement orderby on firestore
func (class *QueryFirestore) OrderBy(path string) Query {
	class.query = class.query.OrderBy(path, firestore.Asc)
	return class
}

//OrderByDesc implement orderby desc on firestore
func (class *QueryFirestore) OrderByDesc(path string) Query {
	class.query = class.query.OrderBy(path, firestore.Desc)
	return class
}

//Limit implement limit on firestore
func (class *QueryFirestore) Limit(n int) Query {
	class.limit = n
	class.query = class.query.Limit(n)
	return class
}

//Offset implement start at on firestore, often use by paginate data
//in firestore will bill extra mony on offset
//func (q *QueryFirestore) Offset(n int) IQuery {
//	q.query = q.query.Offset(n)
//	return q
//}

//Run query with default limit 100 object, use Limit() to override default limit
func (class *QueryFirestore) Run(callback func(o Object)) error {

	if class.limit == 0 {
		class.Limit(100)
	}

	iter := class.query.Documents(class.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		obj := class.newObject()
		err = doc.DataTo(obj)
		if err != nil {
			return err
		}
		obj.SetID(doc.Ref.ID)
		callback(obj)
	}
	return nil
}
