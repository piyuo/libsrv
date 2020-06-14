package data

import "context"

// Query is query interface
type Query interface {
	// Where set where filter
	//
	//	db.Select(ctx, GreetFactory).Where("From", "==", "1").Run(func(o Object) {
	//		i++
	//		err := db.Delete(ctx, o)
	//	})
	//
	Where(path, op string, value interface{}) Query

	// OrderBy set query order by
	//
	//	list = []*Greet{}
	// 	db.Select(ctx, GreetFactory).OrderBy("From").Run(func(o Object) {
	//		greet := o.(*Greet)
	//		list = append(list, greet)
	//	})
	//
	OrderBy(path string) Query

	// OrderByDesc set query order by desc
	//
	//	list = []*Greet{}
	// 	db.Select(ctx, GreetFactory).OrderByDesc("From").Run(func(o Object) {
	//		greet := o.(*Greet)
	//		list = append(list, greet)
	//	})
	//
	OrderByDesc(path string) Query

	// Limit set query limit
	//
	//	list = []*Greet{}
	//	db.Select(ctx, GreetFactory).Limit(1).Run(func(o Object) {
	//		greet := o.(*Greet)
	//		list = append(list, greet)
	//	})
	//
	Limit(n int) Query

	//Offset(n int) IQuery //in firestore will bill extra mony on offset

	// Run query
	//
	//	list = []*Greet{}
	//	db.Select(ctx, GreetFactory).Run(func(o Object) {
	//		greet := o.(*Greet)
	//		list = append(list, greet)
	//	})
	//
	Run(callback func(o Object)) error
}

// AbstractQuery is query object need to implement
type AbstractQuery struct {
	Query
	ctx     context.Context
	factory func() Object
	limit   int
}
