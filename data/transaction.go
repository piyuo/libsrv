package data

import "context"

//Transaction is query interface
type Transaction interface {
	Get(ctx context.Context, obj Object) error
	Put(ctx context.Context, obj Object) error
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
