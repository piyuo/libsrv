package data

import "context"

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
	//	list, err := table.Query().Where("ID", "==", "sample1").Execute(ctx)
	//	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")
	//
	Where(path, op string, value interface{}) Query

	// OrderBy set query order by asc
	//
	//	list, err = table.Query().OrderBy("Name").Execute(ctx)
	//	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")
	//
	OrderBy(path string) Query

	// OrderByDesc set query order by desc
	//
	//	list, err = table.Query().OrderByDesc("Name").Limit(1).Execute(ctx)
	//	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")
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
	//	list, err = table.Query().OrderBy("Name").StartAt("irvine city").Execute(ctx)
	//	So(len(list), ShouldEqual, 1)
	//	So((list[0].(*Sample)).Name, ShouldEqual, "irvine city")
	//
	StartAt(docSnapshotOrFieldValues ...interface{}) Query

	// StartAfter implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
	//
	//	list, err = table.Query().OrderBy("Name").StartAfter("santa ana city").Execute(ctx)
	//	So((list[0].(*Sample)).Name, ShouldEqual, "irvine city")
	//
	StartAfter(docSnapshotOrFieldValues ...interface{}) Query

	// EndAt implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
	//
	//	list, err = table.Query().OrderBy("Name").EndAt("irvine city").Execute(ctx)
	//	So((list[0].(*Sample)).Name, ShouldEqual, "irvine city")
	//
	EndAt(docSnapshotOrFieldValues ...interface{}) Query

	// EndBefore implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
	//
	//	list, err = table.Query().OrderBy("Name").EndBefore("irvine city").Execute(ctx)
	//	So((list[0].(*Sample)).Name, ShouldEqual, "santa ana city")
	//
	EndBefore(docSnapshotOrFieldValues ...interface{}) Query

	// Execute query with default limit to 10 object, use Limit() to override default limit, return nil if anything wrong
	//
	//	list, err = table.Query().OrderByDesc("Name").Limit(1).Execute(ctx)
	//	So(len(list), ShouldEqual, 1)
	//	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")
	//
	Execute(ctx context.Context) ([]Object, error)

	// GetIDs Execute query with default limit to 10 object, use Limit() to override default limit, return nil if anything wrong
	//
	//	idList, err = table.Query().OrderByDesc("Name").Limit(1).GetIDs(ctx)
	//	So(len(idList), ShouldEqual, 1)
	//	So((idList[0], ShouldEqual, "sample2")
	//
	GetIDs(ctx context.Context) ([]string, error)

	// GetFirstObject execute query return first object in result
	//
	//	obj, err := db.Select(ctx, GreetFactory).OrderBy("From").Limit(1).StartAt("b city").GetFirstObject(ctx)
	//	greet := obj.(*Greet)
	//	So(greet.From, ShouldEqual, "b city")
	//
	GetFirstObject(ctx context.Context) (Object, error)

	// GetFirstID execute query return first object id in result
	//
	//	id, err := db.Select(ctx, GreetFactory).OrderBy("From").Limit(1).StartAt("b city").GetFirstID(ctx)
	//	So(id, ShouldEqual, "city1")
	//
	GetFirstID(ctx context.Context) (string, error)

	// Count execute query and return max 10 count
	//
	//	count, err := table.Query().Where("Name", "==", "sample1").Count(ctx)
	//	So(count, ShouldEqual, 1)
	//
	Count(ctx context.Context) (int, error)

	// IsEmpty execute query and return false if object exist
	//
	//	isEmpty, err := table.Query().Where("Name", "==", "sample1").IsEmpty(ctx)
	//	So(isEmpty, ShouldBeFalse)
	//
	IsEmpty(ctx context.Context) (bool, error)

	// IsExists execute query and return true if object exist
	//
	//	isExist, err := table.Query().Where("Name", "==", "sample1").IsExists(ctx)
	//	So(isExist, ShouldBeFalse)
	//
	IsExists(ctx context.Context) (bool, error)

	// Clear keep delete all object in a query until ctx timeout or all object deleted. it delete 500 documents at a time
	//
	Clear(ctx context.Context) (int, error)
}

// BaseQuery represent a query in document database
//
type BaseQuery struct {
	Query

	// Factor use to create object
	//
	Factory func() Object

	// LengthLimit remember if query set limit, if not we will give default limit (10) to avoid return too may document
	//
	LengthLimit int
}
