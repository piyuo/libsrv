package data

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// DBFirestore implement db on firestore
type DBFirestore struct {
	AbstractDB
	client *firestore.Client
}

// Close db connection
//
//	db.Close()
//
func (db *DBFirestore) Close() {
	if db.client != nil {
		db.client.Close()
		db.client = nil
	}
}

// Get data object from data store, return ErrObjectNotFound if object not exist
//
//	err = db.Get(ctx, &greet)
//
func (db *DBFirestore) Get(ctx context.Context, object Object) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	id := object.ID()
	if id == "" {
		return errors.New("get() need object have ID")
	}
	snapshot, err := db.client.Collection(object.ModelName()).Doc(id).Get(ctx)
	if snapshot != nil && !snapshot.Exists() {
		return ErrObjectNotFound
	}
	if err != nil {
		return err
	}

	if err := snapshot.DataTo(object); err != nil {
		return err
	}
	return nil
}

//GetAll object from data store, return error
//
//	err = db.GetAll(ctx, GreetFactory, func(o Object) {}, 100)
//
func (db *DBFirestore) GetAll(ctx context.Context, factory func() Object, limit int, callback func(o Object)) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if limit > 100 {
		panic("GetAll() limit need under 100")
	}
	obj := factory()
	ref := db.client.Collection(obj.ModelName())
	iter := ref.Limit(limit).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return err
		}
		obj := factory()
		err = doc.DataTo(obj)
		if err != nil {
			return err
		}
		obj.SetID(doc.Ref.ID)
		callback(obj)
	}
	return nil
}

// Put data object into data store
//
//	greet := Greet{
//			From:        "1",
//				Description: "1",
//	}
//	ctx := context.Background()
//	db, _ := firestoreNewDB(ctx)
//	defer db.Close()
//	db.Put(ctx, &greet)
func (db *DBFirestore) Put(ctx context.Context, obj Object) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	modelName := obj.ModelName()
	if obj.ID() == "" {
		docRef := db.client.Collection(modelName).NewDoc()
		obj.SetID(docRef.ID)
	}
	_, err := db.client.Collection(modelName).Doc(obj.ID()).Set(ctx, obj)
	if err != nil {
		return errors.Wrap(err, "failed to put object")
	}
	return nil
}

// Update partial object field,  this function is significant slow than put
//
//	err = db.Update(ctx, greet.ModelName(), greet.ID(), map[string]interface{}{
//		"Description": "helloworld",
//	})
//
func (db *DBFirestore) Update(ctx context.Context, modelName string, modelFields map[string]interface{}, objectID string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	_, err := db.client.Collection(modelName).Doc(objectID).Set(ctx, modelFields, firestore.MergeAll)
	if err != nil {
		return errors.Wrap(err, "failed to update field")
	}
	return nil
}

//ListAll get object list from data store, return error
//
//	list, err := db.ListAll(ctx, GreetFactory, 100)
//
func (db *DBFirestore) ListAll(ctx context.Context, factory func() Object, limit int) ([]Object, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if limit > 100 {
		panic("ListAll() limit need under 100")
	}
	obj := factory()
	ref := db.client.Collection(obj.ModelName())
	iter := ref.Limit(limit).Documents(ctx)
	list := []Object{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return list, err
		}
		obj := factory()
		err = doc.DataTo(obj)
		if err != nil {
			return list, err
		}
		obj.SetID(doc.Ref.ID)
		list = append(list, obj)
	}
	return list, nil
}

// Delete object from data store
//
//	_ = db.Delete(ctx, &greet)
//
func (db *DBFirestore) Delete(ctx context.Context, obj Object) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	id := obj.ID()
	modelName := obj.ModelName()
	return db.DeleteByID(ctx, modelName, id)
}

// DeleteByID delete object by id
//
//	err = db.DeleteByID(ctx, GreetModelName, greet.ID())
//
func (db *DBFirestore) DeleteByID(ctx context.Context, modelName, id string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	ref := db.client.Collection(modelName).Doc(id)
	if _, err := ref.Delete(ctx); err != nil {
		return errors.Wrap(err, "failed to delete "+modelName+",id:"+id)
	}
	return nil
}

// DeleteAll delete all object in specific time, return ErrOperationTimeout when timed out
//
//	db.DeleteAll(ctx, GreetModelName, 9)
//
func (db *DBFirestore) DeleteAll(ctx context.Context, modelName string, timeout int) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	beginTime := time.Now()
	ref := db.client.Collection(modelName)
	totalDeleted := 0
	for {
		iter := ref.Limit(100).Documents(ctx)
		numDeleted := 0

		// Iterate through the documents, adding a delete operation for each one to a WriteBatch.
		batch := db.client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return totalDeleted, err
			}
			batch.Delete(doc.Ref)
			numDeleted++
		}

		// If there are no documents to delete  the process is over.
		if numDeleted == 0 {
			return totalDeleted, nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return totalDeleted, err
		}
		totalDeleted += numDeleted
		diff := time.Now().Sub(beginTime).Seconds()
		if int(diff) >= timeout {
			return totalDeleted, ErrOperationTimeout
		}
	}
}

// Select create query
//
//	query := db.Select(ctx, func() Object {
//		return new(Greet)
//	})
//
func (db *DBFirestore) Select(ctx context.Context, factory func() Object) Query {
	if factory == nil {
		panic("Select must have factory function like func(){new(object)}")
	}
	obj := factory()
	query := db.client.Collection(obj.ModelName()).Query
	return NewQueryFirestore(ctx, query, factory)
}

// Transaction start a transaction
//
//	err := db.Transaction(ctx, func(ctx context.Context, tx Transaction) error {
//		tx.Put(ctx, &greet1)
//		tx.Put(ctx, &greet2)
//		return nil
//	})
//
func (db *DBFirestore) Transaction(ctx context.Context, callback func(ctx context.Context, tx Transaction) error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return db.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		tf := NewTransactionFirestore(ctx, db.client, tx)
		return callback(ctx, tf)
	})
}

// Exist return true if query result return a least one object
//
//	exist, err := db.Exist(ctx, GreetModelName, "From", "==", "1")
//
func (db *DBFirestore) Exist(ctx context.Context, modelName, modelField, operator string, value interface{}) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	docIterator := db.client.Collection(modelName).Query.Where(modelField, operator, value).Limit(1).Documents(ctx)
	defer docIterator.Stop()

	_, err := docIterator.Next()
	if err == iterator.Done {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// ExistByID return true if query result return a least one object
//
//	exist, err := db.ExistByID(ctx, GreetModelName, "greet1")
//
func (db *DBFirestore) ExistByID(ctx context.Context, modelName, id string) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	snapshot, err := db.client.Collection(modelName).Doc(id).Get(ctx)
	if snapshot != nil && !snapshot.Exists() {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "failed to find exist "+modelName+", id:"+id)
	}
	return true, nil
}

// Count10 return max 10 result set,cause firestore are charged for a read each time a document in the result set, we need keep result set as small as possible
//
//	count, err := db.Count10(ctx, GreetModelName, "From", "==", "1")
//
func (db *DBFirestore) Count10(ctx context.Context, modelName, modelField, operator string, value interface{}) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	docIterator := db.client.Collection(modelName).Query.Where(modelField, operator, value).Limit(10).Documents(ctx)
	defer docIterator.Stop()
	count := 0
	for {
		_, err := docIterator.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return 0, err
		}
		count++
	}
	return count, nil
}

// Increment value on object field
//
//	err := db.Increment(ctx, GreetModelName, "Value", greet.ID(), 2)
//
func (db *DBFirestore) Increment(ctx context.Context, modelName, modelField, objectID string, value int) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	doc := db.client.Collection(modelName).Doc(objectID)

	//Update() return writeResult, we don't need
	_, err := doc.Update(ctx, []firestore.Update{
		{Path: modelField, Value: firestore.Increment(value)},
	})
	if err != nil {
		return err
	}
	return nil
}

// Counter get counter from data store, create one if not exist
//
//	counter,err = db.Counter(ctx, "myCounter",10)
//
func (db *DBFirestore) Counter(ctx context.Context, name string, numShards int) (*Counter, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if numShards <= 0 {
		numShards = 10
	}
	counter := &Counter{
		NumShards: numShards,
	}
	counter.docRef = db.client.Collection("Counter").Doc(name)
	snapshot, err := counter.docRef.Get(ctx)

	if snapshot != nil && !snapshot.Exists() {
		_, err := counter.docRef.Set(ctx, counter)
		if err != nil {
			return nil, errors.Wrap(err, "failed to put counter")
		}
		err = counter.init(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init counter")
		}
		return counter, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to get counter")
	}

	if err := snapshot.DataTo(counter); err != nil {
		return nil, errors.Wrap(err, "failed convert to counter")
	}
	return counter, nil
}

// DeleteCounter remove counter and all shards
//
//	counter,err = db.GetCounter(ctx, "myCounter")
//
func (db *DBFirestore) DeleteCounter(ctx context.Context, name string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	shardsBatch := db.client.Batch()
	docRef := db.client.Collection("Counter").Doc(name)
	shards := docRef.Collection("shards").Documents(ctx)
	numDeleted := 0
	for {
		doc, err := shards.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to delete shards")
		}
		shardsBatch.Delete(doc.Ref)
		numDeleted++
	}

	if numDeleted > 0 {
		_, err := shardsBatch.Commit(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to commit delete shards")
		}

	}

	if _, err := docRef.Delete(ctx); err != nil {
		return errors.Wrap(err, "failed to delete counter:"+name)
	}

	return nil
}
