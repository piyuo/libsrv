package data

import "context"

//Transaction is query interface
type Transaction interface {
	Get(ctx context.Context, obj Object) error
	Put(ctx context.Context, obj Object) error
	Delete(ctx context.Context, obj Object) error
}
