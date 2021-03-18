package db

import (
	"context"

	"cloud.google.com/go/firestore"
)

// BatchFunc define a batch function
//
type BatchFunc func(ctx context.Context, bc Batch) error

// Batch define batch operation
//
type Batch interface {

	// Set object into table, If the document not exist, it will be created. If the document does exist, its contents will be overwritten with the newly provided data, if object does not have id, it will created using UUID
	//
	//	 err := Set(ctx, object)
	//
	Set(ctx context.Context, obj Object) error

	// Update partial object field, create new one if object does not exist, this function is significant slow than Set()
	//
	//	err = Update(ctx, Sample, map[string]interface{}{
	//		"desc": "hi",
	//	})
	//
	Update(ctx context.Context, obj Object, fields map[string]interface{}) error

	// Increment value on object field, return error if object does not exist
	//
	//	err := Increment(ctx,sample, "Value", 2)
	//
	Increment(ctx context.Context, obj Object, field string, value int) error

	// Delete object, no error if id not exist
	//
	//	Delete(ctx, sample)
	//
	Delete(ctx context.Context, obj Object) error

	// DeleteList delete object use list of id, no error if id not exist
	//
	//	c.DeleteList(ctx, &Sample{}, []string{"1","2"})
	//
	DeleteList(ctx context.Context, obj Object, list []string) error

	// DeleteRef delete object use document ref
	//
	//	err := DeleteRef(ref)
	//
	DeleteRef(ref *firestore.DocumentRef)
}
