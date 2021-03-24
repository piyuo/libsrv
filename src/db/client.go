package db

import (
	"context"
	"time"

	"github.com/piyuo/libsrv/src/env"
	"github.com/piyuo/libsrv/src/identifier"
	"github.com/pkg/errors"
)

// Client define how to connect and manipulate document database
//
type Client interface {
	// Close database connection
	//
	//	Close()
	//
	Close()

	// IsClose return true if connection is close
	//
	//	closed := IsClose()
	//
	IsClose() bool

	// Get data object from data store, return nil if object does not exist
	//
	//	object, err := Get(ctx, &Sample{}, "id")
	//
	Get(ctx context.Context, obj Object, id string) (Object, error)

	// Exists return true if object with id exist
	//
	//	found,err := Exists(ctx, &Sample{}, "id")
	//
	Exists(ctx context.Context, obj Object, id string) (bool, error)

	// List return object list, use max to specific return object count
	//
	//	list,err := List(ctx, &Sample{},10)
	//
	List(ctx context.Context, obj Object, max int) ([]Object, error)

	// Select return object field from data store, return nil if object does not exist
	//
	//	return Select(ctx, &Sample{}, id, field)
	//
	Select(ctx context.Context, obj Object, id, field string) (interface{}, error)

	// Query create query
	//
	//	c.Query(ctx, &Sample{}).Execute(ctx)
	//
	Query(obj Object) Query

	// Set object into data store, If the document does not exist, it will be created. If the document does exist, its contents will be overwritten with the newly provided data,
	// if object does not have id, it will created using UUID
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

	// DeleteAll delete all document in collection
	//
	//	done, err := Truncate(ctx, "Sample", 50)
	//
	Truncate(ctx context.Context, collectionName string, max int) (bool, error)

	// Transaction start a transaction operation
	//
	//	err := Transaction(ctx, func(ctx context.Context,tx db.Transaction) error {
	//		return nil
	//	})
	//
	Transaction(ctx context.Context, f TransactionFunc) error

	// Batch start a batch operation. batch won't be commit if there is no batch operation like set/update/delete been called
	//
	//	err := Batch(ctx, func(ctx context.Context,bc db.Batch) error {
	//		return nil
	//	})
	//
	Batch(ctx context.Context, f BatchFunc) error

	// Counter return counter, set numshards 100 times of concurrent usage. for example if you think concurrent use is 10/seconds then set numshards to 1000 to avoid too much retention error
	// if DateHierarchyFull is set, counter will automatically generate year/month/day/hour hierarchy in utc timezone
	//
	//	sampleCounter,err = Counter("SampleCount", 100, DateHierarchyNone)
	//
	Counter(counterName string, numshards int, hierarchy DateHierarchy) Counter

	// Coder return coder, set numshards 100 times of concurrent usage. for example if you think concurrent use is 10/seconds then set numshards to 1000 to avoid too much retention error
	//
	//	productCoder,err = Coder("coderName",100)
	//
	Coder(coderName string, numshards int) Coder

	// Serial return serial, create one if not exist, please be aware Serial can only generate 1 number per second, use serial with high frequency will cause too much retention error
	//
	//	productNo,err = Serial("serialName")
	//
	Serial(serialName string) Serial
}

// DateHierarchy used in create counter
//
type DateHierarchy int8

const (
	// DateHierarchyNone create counter without date hierarchy, only total count
	//
	DateHierarchyNone DateHierarchy = 1

	// DateHierarchyFull create counter with year/month/day/hour hierarchy and total count
	//
	DateHierarchyFull = 2
)

// AssertObject return error if ctx/obj has problem
//
//	err := AssertObject(ctx, Sample)
//
func AssertObject(ctx context.Context, obj Object, hasID bool) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if obj == nil {
		return errors.New("obj must not nil")
	}
	if hasID && obj.ID() == "" {
		return errors.New(obj.Collection() + " obj must has ID")
	}
	return nil
}

// AssertID check id is empty
//
//	err := AssertID(id)
//
func AssertID(id string) error {
	if id == "" {
		return errors.New("id must not empty")
	}
	return nil
}

// BaseClient represent object stored in document database
//
type BaseClient struct {
	Client
}

// BeforeSet add create/update time, accountID, userID before set to database
//
func (c *BaseClient) BeforeSet(ctx context.Context, obj Object) {
	if obj.ID() == "" {
		obj.SetID(identifier.UUID())
	}

	t := time.Now().UTC()
	obj.SetCreateTime(t)
	obj.SetUpdateTime(t)

	accountID := env.GetAccountID(ctx)
	if accountID != "" {
		obj.SetAccountID(accountID)
	}

	userID := env.GetUserID(ctx)
	if userID != "" {
		obj.SetUserID(userID)
	}
}
