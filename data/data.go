package data

import (
	"context"
)

const here = "data"

// NewGlobalDB create global database from default provider
//
//	db := data.NewGlobalDB(ctx)
//	err := db.Put(ctx, &greet)
//	retrive := Greet{}
//	retrive.SetID(greet.ID())
//	err = db.Get(ctx, &retrive)
func NewGlobalDB(ctx context.Context) (DB, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return firestoreGlobalDB(ctx)
}

// NewRegionalDB create regional database from default provider
//
//	db := data.NewGlobalDB(ctx)
//	err := db.Put(ctx, &greet)
//	retrive := Greet{}
//	retrive.SetID(greet.ID())
//	err = db.Get(ctx, &retrive)
func NewRegionalDB(ctx context.Context) (DB, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return firestoreRegionalDB(ctx)
}
