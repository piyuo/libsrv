package data

import "context"

// QueryRef represent query public method
type QueryRef interface {
	// Where set where filter
	//
	//	list, err := table.Query().Where("Name", "==", "sample1").Execute(ctx)
	//	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")
	//
	Where(path, op string, value interface{}) QueryRef

	// OrderBy set query order by asc
	//
	//	list, err = table.Query().OrderBy("Name").Execute(ctx)
	//	So((list[0].(*Sample)).Name, ShouldEqual, "sample1")
	//
	OrderBy(path string) QueryRef

	// OrderByDesc set query order by desc
	//
	//	list, err = table.Query().OrderByDesc("Name").Limit(1).Execute(ctx)
	//	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")
	//
	OrderByDesc(path string) QueryRef

	// Limit set query limit
	//
	//	list, err = table.Query().OrderBy("Name").Limit(1).Execute(ctx)
	//	So(len(list), ShouldEqual, 1)
	//
	Limit(n int) QueryRef

	// StartAt implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
	//
	//	list, err = table.Query().OrderBy("Name").StartAt("irvine city").Execute(ctx)
	//	So(len(list), ShouldEqual, 1)
	//	So((list[0].(*Sample)).Name, ShouldEqual, "irvine city")
	//
	StartAt(docSnapshotOrFieldValues ...interface{}) QueryRef

	// StartAfter implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
	//
	//	list, err = table.Query().OrderBy("Name").StartAfter("santa ana city").Execute(ctx)
	//	So((list[0].(*Sample)).Name, ShouldEqual, "irvine city")
	//
	StartAfter(docSnapshotOrFieldValues ...interface{}) QueryRef

	// EndAt implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
	//
	//	list, err = table.Query().OrderBy("Name").EndAt("irvine city").Execute(ctx)
	//	So((list[0].(*Sample)).Name, ShouldEqual, "irvine city")

	//
	EndAt(docSnapshotOrFieldValues ...interface{}) QueryRef

	// EndBefore implement Paginate on firestore, please be aware not use index but fieldValue to do the trick, see sample
	//
	//	list, err = table.Query().OrderBy("Name").EndBefore("irvine city").Execute(ctx)
	//	So((list[0].(*Sample)).Name, ShouldEqual, "santa ana city")
	//
	EndBefore(docSnapshotOrFieldValues ...interface{}) QueryRef

	// Execute query with default limit to 10 object, use Limit() to override default limit, return nil if anything wrong
	//
	//	list, err = table.Query().OrderByDesc("Name").Limit(1).Execute(ctx)
	//	So(len(list), ShouldEqual, 1)
	//	So((list[0].(*Sample)).Name, ShouldEqual, "sample2")
	//
	Execute(ctx context.Context) ([]ObjectRef, error)

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
}

// Query represent a query in document database
//
type Query struct {
	QueryRef

	// factor use to create object
	//
	factory func() ObjectRef

	// limit remember if query set limit, if not we will give default limit (10) to avoid return too may document
	//
	limit int
}
