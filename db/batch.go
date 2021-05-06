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
	//	 Set(ctx, object) // no error in batch mode
	//
	Set(ctx context.Context, obj Object)

	// Update partial object field, create new one if object does not exist, this function is significant slow than Set()
	//
	//	Update(ctx, Sample, map[string]interface{}{
	//		"desc": "hi",
	//	})
	//
	Update(ctx context.Context, obj Object, fields map[string]interface{})

	// Increment value on object field, return error if object does not exist
	//
	//	Increment(ctx,sample, "Value", 2) // no error in batch mode
	//
	Increment(ctx context.Context, obj Object, field string, value int)

	// Delete object, no error if id not exist
	//
	//	Delete(ctx, sample) // no error in batch mode
	//
	Delete(ctx context.Context, obj Object)

	// DeleteList delete object use list of id, no error if id not exist
	//
	//	DeleteList(ctx, &Sample{}, []string{"1","2"}) // no error in batch mode
	//
	DeleteList(ctx context.Context, obj Object, list []string)

	// DeleteRef delete object use document ref
	//
	//	DeleteRef(ref) // no error in batch mode
	//
	DeleteRef(ref *firestore.DocumentRef)
}
