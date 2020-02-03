package data

import (
	"context"
)

const here = "data"

// NewDB create db from default provider
//
//	db := data.NewDB(ctx)
//	err := db.Put(ctx, &greet)
//	retrive := Greet{}
//	retrive.SetID(greet.ID())
//	err = db.Get(ctx, &retrive)
func NewDB(ctx context.Context) (DB, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return firestoreNewDB(ctx)
}
