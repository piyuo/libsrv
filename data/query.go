package data

// IQuery is query interface
type IQuery interface {
	Where(path, op string, value interface{}) IQuery
	OrderBy(path string) IQuery
	OrderByDesc(path string) IQuery
	Limit(n int) IQuery
	//Offset(n int) IQuery //in firestore will bill extra mony on offset
	Run(callback func(o IObject)) error
}

// Query represent db query
type Query struct {
}
