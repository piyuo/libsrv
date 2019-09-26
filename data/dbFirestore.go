package data

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// DBFirestore implement db on firestore
type DBFirestore struct {
	DB
	client *firestore.Client
	ctx    context.Context
}

//NewDBFirestore new firestore db instance
func NewDBFirestore(client *firestore.Client) *DBFirestore {
	db := &DBFirestore{}
	db.client = client
	db.ctx = context.Background()
	return db
}

//Close a db connection
func (db *DBFirestore) Close() {
	if db.client != nil {
		db.client.Close()
		db.client = nil
	}
}

//Put data object into data store
func (db *DBFirestore) Put(obj IObject) error {
	Class := obj.Class()
	if obj.ID() == "" {
		ref := db.client.Collection(Class).NewDoc()
		obj.SetID(ref.ID)
	}
	_, err := db.client.Collection(Class).Doc(obj.ID()).Set(db.ctx, obj)
	if err != nil {
		return errors.Wrap(err, "failed to put object")
	}
	return nil
}

//Update partial object field  in firestore,  this function is not significant fast than put
func (db *DBFirestore) Update(objClass string, objID string, fields map[string]interface{}) error {
	_, err := db.client.Collection(objClass).Doc(objID).Set(db.ctx, fields, firestore.MergeAll)
	if err != nil {
		return errors.Wrap(err, "failed to update field")
	}
	return nil
}

//Get data object from data store, return ErrNotFound if object not exist
func (db *DBFirestore) Get(obj IObject) error {
	return db.GetByClass(obj.Class(), obj)
}

//GetByClass get object from data store,use class instead of obj class
func (db *DBFirestore) GetByClass(class string, obj IObject) error {
	id := obj.ID()
	if id == "" {
		return ErrNotFound
	}
	snapshot, err := db.client.Collection(class).Doc(id).Get(db.ctx)
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return ErrNotFound
		}
		return err
	}
	if err := snapshot.DataTo(obj); err != nil {
		return err
	}
	return nil
}

//GetAll object from data store, return error
func (db *DBFirestore) GetAll(factory func() IObject, callback func(o IObject), limit int) error {
	if limit > 100 {
		panic("GetAll() limit need under 100")
	}
	obj := factory()
	ref := db.client.Collection(obj.Class())
	iter := ref.Limit(limit).Documents(db.ctx)

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

//ListAll get object lit from data store, return error
func (db *DBFirestore) ListAll(factory func() IObject, limit int) ([]IObject, error) {
	if limit > 100 {
		panic("ListAll() limit need under 100")
	}
	obj := factory()
	ref := db.client.Collection(obj.Class())
	iter := ref.Limit(limit).Documents(db.ctx)
	list := []IObject{}
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
func (db *DBFirestore) Delete(obj IObject) error {
	id := obj.ID()
	class := obj.Class()
	ref := db.client.Collection(class).Doc(id)
	if _, err := ref.Delete(db.ctx); err != nil {
		return err
	}
	return nil
}

//DeleteAll only run in time, return ErrTimeout when time is up
func (db *DBFirestore) DeleteAll(className string, timeout int) (int, error) {
	beginTime := time.Now()
	ref := db.client.Collection(className)
	totalDeleted := 0
	for {
		iter := ref.Limit(100).Documents(db.ctx)
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

		_, err := batch.Commit(db.ctx)
		if err != nil {
			return totalDeleted, err
		}
		totalDeleted += numDeleted
		diff := time.Now().Sub(beginTime).Seconds()
		if int(diff) >= timeout {
			return totalDeleted, ErrTimeout
		}
	}
}

//Select data object from firestore
func (db *DBFirestore) Select(f func() IObject) IQuery {
	if f == nil {
		panic("Select must have new function like func(){new(object)}")
	}
	obj := f()
	query := db.client.Collection(obj.Class()).Query
	return NewQueryFirestore(db.ctx, query, f)
}

//RunTransaction implement firestore run transaction
func (db *DBFirestore) RunTransaction(f func(tx ITransaction) error) error {
	return db.client.RunTransaction(db.ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		tf := NewTransactionFirestore(db.ctx, db.client, tx)
		return f(tf)
	})
}
