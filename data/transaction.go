package data

import "context"

//Transaction is query interface
type Transaction interface {
	// Get data object from data store, return ErrNotFound if object not exist
	//
	//	greet := &Greet{}
	//	greet.SetID(greet1.ID())
	//	err = db.Transaction(ctx, func(ctx context.Context, tx Transaction) error {
	//		tx.Get(ctx, greet)
	//		return nil
	//	})
	//
	Get(ctx context.Context, obj Object) error

	// Put data object into data store
	//
	//	greet1 := &Greet{
	//		From:        "1",
	//		Description: "1",
	//	}
	//	err := db.Transaction(ctx, func(ctx context.Context, tx Transaction) error {
	//		tx.Put(ctx, greet1)
	//		return nil
	//	})
	Put(ctx context.Context, obj Object) error

	// Delete data object from firestore
	//
	//	err := db.Transaction(ctx, func(ctx context.Context, tx Transaction) error {
	//		tx.Delete(ctx, greet)
	//		return nil
	//	})
	//
	Delete(ctx context.Context, obj Object) error

	// ShortID create unique serial number, please be aware serial can only generate one number per second and use with transation to ensure unique
	//
	//	id,err = db.ShortID(ctx, "myID")
	//
	ShortID(ctx context.Context, name string) (*ShortID, error)
}

// AbstractTransaction is parent class for all DB child
type AbstractTransaction struct {
	Transaction
}
