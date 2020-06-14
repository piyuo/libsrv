package data

import (
	"context"

	"github.com/pkg/errors"
)

// DB represent database
type DB interface {
	// Close db connection
	//
	//	db.Close()
	//
	Close()

	// Put data object into data store
	//
	//	greet := &Greet{
	//		From:        "1",
	//		Description: "1",
	//	}
	//	ctx := context.Background()
	//	db, _ := firestoreNewDB(ctx)
	//	defer db.Close()
	//	db.Put(ctx, greet)
	Put(ctx context.Context, obj Object) error

	// Update partial object field,  this function is significant slow than put
	//
	//	ctx := context.Background()
	//	db, _ := firestoreNewDB(ctx)
	//	defer db.Close()
	//	err = db.Update(ctx, greet.ModelName(), greet.ID(), map[string]interface{}{
	//		"Description": "helloworld",
	//	})
	//
	Update(ctx context.Context, modelName string, modelFields map[string]interface{}, objectID string) error

	// Get data object from data store, return ErrObjectNotFound if object not exist
	//
	//	ctx := context.Background()
	//	db, _ := firestoreNewDB(ctx)
	//	defer db.Close()
	//	greet := &Greet{}
	//	greet.SetID("myID")
	//	err = db.Get(ctx, &greet)
	//
	Get(ctx context.Context, obj Object) error

	//GetAll object from data store, return error
	//
	//	err = db.GetAll(ctx, GreetFactory, func(o Object) {}, 100)
	//
	GetAll(ctx context.Context, factory func() Object, limit int, callback func(o Object)) error

	//ListAll get object list from data store, return error
	//
	//	list, err := db.ListAll(ctx, GreetFactory, 100)
	//
	ListAll(ctx context.Context, factory func() Object, limit int) ([]Object, error)

	// Delete object from data store
	//
	//	ctx := context.Background()
	//	db, _ := firestoreNewDB(ctx)
	//	defer db.Close()
	//	greet := &Greet{}
	//	db.Put(ctx, &greet)
	//	_ = db.Delete(ctx, greet)
	//
	Delete(ctx context.Context, obj Object) error

	// Delete object from data store using model name and id
	//
	//	greet := Greet{
	//		From:        "me",
	//		Description: "hi",
	//	}
	//	ctx := context.Background()
	//	db, err := firestoreGlobalDB(ctx)
	//	defer db.Close()
	//	db.Put(ctx, &greet)
	//	err = db.DeleteByID(ctx, GreetModelName, greet.ID())
	//
	DeleteByID(ctx context.Context, modelName, id string) error

	// DeleteAll delete all object in specific time, return ErrOperationTimeout when timed out
	//
	//	db.DeleteAll(ctx, GreetModelName, 9)
	//
	DeleteAll(ctx context.Context, modelName string, timeout int) (int, error)

	// Select create query
	//
	//	query := db.Select(ctx, func() Object {
	//		return new(Greet)
	//	})
	//
	Select(ctx context.Context, factory func() Object) Query

	// Transaction start a transaction
	//
	//	err := db.Transaction(ctx, func(ctx context.Context, tx Transaction) error {
	//		tx.Put(ctx, &greet1)
	//		return nil
	//	})
	//
	Transaction(ctx context.Context, callback func(ctx context.Context, tx Transaction) error) error

	// Exist return true if query result return a least one object
	//
	//	exist, err := db.Exist(ctx, GreetModelName, "From", "==", "1")
	//
	Exist(ctx context.Context, modelName, modelField, operator string, value interface{}) (bool, error)

	// ExistByID return true if query result return a least one object
	//
	//	exist, err := db.ExistByID(ctx, GreetModelName, "greet1")
	//
	ExistByID(ctx context.Context, modelName, id string) (bool, error)

	// Count10 return max 10 result set,cause firestore are charged for a read each time a document in the result set, we need keep result set as small as possible
	//
	//	count, err := db.Count10(ctx, GreetModelName, "From", "==", "1")
	//
	Count10(ctx context.Context, modelName, modelField, operator string, value interface{}) (int, error)

	// Increment value on object field
	//
	//	err := db.Increment(ctx, GreetModelName, "Value", greet.ID(), 2)
	//
	Increment(ctx context.Context, modelName, modelField, objectID string, value int) error

	// Get counter from data store, create one if not exist
	//
	//	counter,err = db.Counter(ctx, "myCounter",10)
	//
	Counter(ctx context.Context, name string, numShards int) (*Counter, error)

	// DeleteCounter delete remove counter and all shards
	//
	//	err = db.DeleteCounter(ctx, "myCounter")
	//
	DeleteCounter(ctx context.Context, name string) error
}

// AbstractDB is parent class for all DB child
type AbstractDB struct {
	DB
}

// ErrOperationTimeout is returned by DeleteAll method when the method is run too long
var ErrOperationTimeout = errors.New("db operation timeout")

// ErrObjectNotFound is returned by Get method object not exist
var ErrObjectNotFound = errors.New("object not found")
