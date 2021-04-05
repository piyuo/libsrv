package db

import (
	"context"
)

// orderby define order by asc or desc
//
type orderby int

const (
	// ASC mean order by asc
	//
	ASC orderby = iota

	// DESC mean order by DESC
	//
	DESC
)

// Query represent query public method
type Query interface {
	// Where set filter, if path == "ID" mean using document id in as filter
	//
	//	list, err := Query(&Sample{}).Where("ID", "==", "sample1").Execute(ctx)
	//
	Where(path, op string, value interface{}) Query

	// OrderBy set query order by asc
	//
	//	list, err = Query(&Sample{}).OrderBy("Name").Execute(ctx)
	//
	OrderBy(path string) Query

	// Limit set query limit
	//
	//	list, err = Query(&Sample{}).OrderBy("Name").Limit(1).Execute(ctx)
	//
	OrderByDesc(path string) Query

	// Limit set query limit
	//
	//	list, err = table.Query().OrderBy("Name").Limit(1).Execute(ctx)
	//	So(len(list), ShouldEqual, 1)
	//
	Limit(n int) Query

	// StartAt implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
	//
	//	list, err = Query(&Sample{}).OrderBy("Name").StartAt("irvine city").Execute(ctx)
	//
	StartAt(docSnapshotOrFieldValues ...interface{}) Query

	// StartAfter implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
	//
	//	list, err = Query(&Sample{}).OrderBy("Name").StartAfter("santa ana city").Execute(ctx)
	//
	StartAfter(docSnapshotOrFieldValues ...interface{}) Query

	// EndAt implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
	//
	//	list, err = Query(&Sample{}).OrderBy("Name").EndAt("irvine city").Execute(ctx)
	//
	EndAt(docSnapshotOrFieldValues ...interface{}) Query

	// EndBefore implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
	//
	//	list, err = Query(&Sample{}).OrderBy("Name").EndBefore("irvine city").Execute(ctx)
	//
	EndBefore(docSnapshotOrFieldValues ...interface{}) Query

	// Delete delete all document return from query. delete max doc count. return is done,delete count, error
	//
	//	done, count, err := client.Query(&Sample{}).Where("Name", "==", name).Delete(ctx, 100)
	//
	Delete(ctx context.Context, max int) (bool, int, error)

	// Cleanup delete 25 document a time, max 1000 object. return true if no object left in collection
	//
	//	done,  err := client.Query(&Sample{}).Where("Name", "==", name).Cleanup(ctx)
	//
	Cleanup(ctx context.Context) (bool, error)

	// Return query result with default limit to 20 object, use Limit() to override default limit, return nil if anything wrong
	//
	//	list, err = Query(&Sample{}).OrderByDesc("Name").Limit(1).Return(ctx)
	//
	Return(ctx context.Context) ([]Object, error)

	// ReturnID only return object id with default limit to 20 object, use Limit() to override default limit, return nil if anything wrong
	//
	//	idList, err := Query(&Sample{}).OrderBy("From").Limit(1).StartAt("b city").ReturnID(ctx)
	//
	ReturnID(ctx context.Context) ([]string, error)

	// ReturnCount return object count with default limit to 20 object, use Limit() to override default limit
	//
	//	count, err := Query(&Sample{}).Where("Name", "==", "sample1").ReturnCount(ctx)
	//
	ReturnCount(ctx context.Context) (int, error)

	// ReturnEmpty return true if no object exist
	//
	//	isEmpty, err := Query(&Sample{}).Where("Name", "==", "sample1").ReturnEmpty(ctx)
	//
	ReturnEmpty(ctx context.Context) (bool, error)

	// ReturnExists return true if object exist
	//
	//	isExists, err := Query(&Sample{}).Where("Name", "==", "sample1").ReturnExists(ctx)
	//
	ReturnExists(ctx context.Context) (bool, error)

	// ReturnFirst return first object from query
	//
	//	obj, err := Query(&Sample{}).OrderBy("From").Limit(1).StartAt("b city").ReturnFirst(ctx)
	//	greet := obj.(*Greet)
	//
	ReturnFirst(ctx context.Context) (Object, error)

	// ReturnFirstID return first object id from query
	//
	//	id, err := Query(&Sample{}).OrderBy("From").Limit(1).StartAt("b city").ReturnFirstID(ctx)
	//
	ReturnFirstID(ctx context.Context) (string, error)
}

// BaseQuery represent a query in document database
//
type BaseQuery struct {
	Query

	// QueryObject for query and create object
	//
	QueryObject Object

	// QueryTransaction not nil mean using transaction to do query
	//
	QueryTransaction Transaction
}
