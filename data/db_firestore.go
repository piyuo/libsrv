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
func (db *DBFirestore) Close() {
	if db.client != nil {
		db.client.Close()
		db.client = nil
	}
}

// Get data object from data store, return ErrNotFound if object not exist
//
//	err = db.Get(ctx, &greet)
func (db *DBFirestore) Get(ctx context.Context, obj Object) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return db.GetByModelName(ctx, obj.ModelName(), obj)
}

//GetByModelName get object from data store,use class instead of obj class
//
//	err = db.GetByClass(ctx, "Greet", &greet)
func (db *DBFirestore) GetByModelName(ctx context.Context, modelName string, obj Object) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	id := obj.ID()
	if id == "" {
		return errors.New("get object need object  have ID")
	}
	snapshot, err := db.client.Collection(modelName).Doc(id).Get(ctx)
	if snapshot != nil && !snapshot.Exists() {
		return ErrObjectNotFound
	}
	if err != nil {
		return err
	}

	if err := snapshot.DataTo(obj); err != nil {
		return err
	}
	return nil
}

//GetAll object from data store, return error
func (db *DBFirestore) GetAll(ctx context.Context, factory func() Object, callback func(o Object), limit int) error {
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
	Class := obj.ModelName()
	if obj.ID() == "" {
		ref := db.client.Collection(Class).NewDoc()
		obj.SetID(ref.ID)
	}
	_, err := db.client.Collection(Class).Doc(obj.ID()).Set(ctx, obj)
	if err != nil {
		return errors.Wrap(err, "failed to put object")
	}
	return nil
}

//Update partial object field  in firestore,  this function is not significant fast than put
func (db *DBFirestore) Update(ctx context.Context, objClass string, objID string, fields map[string]interface{}) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	_, err := db.client.Collection(objClass).Doc(objID).Set(ctx, fields, firestore.MergeAll)
	if err != nil {
		return errors.Wrap(err, "failed to update field")
	}
	return nil
}

//ListAll get object lit from data store, return error
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

//Delete data object from data store
func (db *DBFirestore) Delete(ctx context.Context, obj Object) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	id := obj.ID()
	class := obj.ModelName()
	ref := db.client.Collection(class).Doc(id)
	if _, err := ref.Delete(ctx); err != nil {
		return err
	}
	return nil
}

// DeleteAll only run in time, return ErrTimeout when time is up
//
//	db.DeleteAll(ctx, GreetModelName, 9)
//
func (db *DBFirestore) DeleteAll(ctx context.Context, className string, timeout int) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}
	beginTime := time.Now()
	ref := db.client.Collection(className)
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

// Select data object from firestore
func (db *DBFirestore) Select(ctx context.Context, f func() Object) Query {
	if f == nil {
		panic("Select must have new function like func(){new(object)}")
	}
	obj := f()
	query := db.client.Collection(obj.ModelName()).Query
	return NewQueryFirestore(ctx, query, f)
}

// RunTransaction implement firestore run transaction
func (db *DBFirestore) RunTransaction(ctx context.Context, f func(ctx context.Context, tx Transaction) error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return db.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		tf := NewTransactionFirestore(ctx, db.client, tx)
		return f(ctx, tf)
	})
}

// Exist return true if query result return a least one object
//
//	exist, err := db.Exist(ctx, GreetModelName, "From", "==", "1")
//
func (db *DBFirestore) Exist(ctx context.Context, class, path, op string, value interface{}) (bool, error) {
	docIterator := db.client.Collection(class).Query.Where(path, op, value).Limit(1).Documents(ctx)
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

// Count10 return max 10 result set,cause firestore are charged for a read each time a document in the result set, we need keep result set as small as possible
//
//	count, err := db.Count10(ctx, GreetModelName, "From", "==", "1")
//
func (db *DBFirestore) Count10(ctx context.Context, class, path, op string, value interface{}) (int, error) {
	docIterator := db.client.Collection(class).Query.Where(path, op, value).Limit(10).Documents(ctx)
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
