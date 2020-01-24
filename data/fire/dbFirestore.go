package fire

import (
	"context"
	"time"

	"github.com/piyuo/go-libsrv/data/protocol"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// DBFirestore implement db on firestore
type DBFirestore struct {
	protocol.DB
	client *firestore.Client
}

//NewDBFirestore new firestore db instance
func NewDBFirestore(client *firestore.Client) *DBFirestore {
	db := &DBFirestore{}
	db.client = client
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
func (db *DBFirestore) Put(ctx context.Context, obj protocol.Object) error {
	Class := obj.Class()
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
	_, err := db.client.Collection(objClass).Doc(objID).Set(ctx, fields, firestore.MergeAll)
	if err != nil {
		return errors.Wrap(err, "failed to update field")
	}
	return nil
}

//Get data object from data store, return ErrNotFound if object not exist
func (db *DBFirestore) Get(ctx context.Context, obj protocol.Object) error {
	return db.GetByClass(ctx, obj.Class(), obj)
}

//GetByClass get object from data store,use class instead of obj class
func (db *DBFirestore) GetByClass(ctx context.Context, class string, obj protocol.Object) error {
	id := obj.ID()
	if id == "" {
		return protocol.ErrNotFound
	}
	snapshot, err := db.client.Collection(class).Doc(id).Get(ctx)
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return protocol.ErrNotFound
		}
		return err
	}
	if err := snapshot.DataTo(obj); err != nil {
		return err
	}
	return nil
}

//GetAll object from data store, return error
func (db *DBFirestore) GetAll(ctx context.Context, factory func() protocol.Object, callback func(o protocol.Object), limit int) error {
	if limit > 100 {
		panic("GetAll() limit need under 100")
	}
	obj := factory()
	ref := db.client.Collection(obj.Class())
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

//ListAll get object lit from data store, return error
func (db *DBFirestore) ListAll(ctx context.Context, factory func() protocol.Object, limit int) ([]protocol.Object, error) {
	if limit > 100 {
		panic("ListAll() limit need under 100")
	}
	obj := factory()
	ref := db.client.Collection(obj.Class())
	iter := ref.Limit(limit).Documents(ctx)
	list := []protocol.Object{}
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
func (db *DBFirestore) Delete(ctx context.Context, obj protocol.Object) error {
	id := obj.ID()
	class := obj.Class()
	ref := db.client.Collection(class).Doc(id)
	if _, err := ref.Delete(ctx); err != nil {
		return err
	}
	return nil
}

//DeleteAll only run in time, return ErrTimeout when time is up
func (db *DBFirestore) DeleteAll(ctx context.Context, className string, timeout int) (int, error) {
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
			return totalDeleted, protocol.ErrTimeout
		}
	}
}

//Select data object from firestore
func (db *DBFirestore) Select(ctx context.Context, f func() protocol.Object) protocol.Query {
	if f == nil {
		panic("Select must have new function like func(){new(object)}")
	}
	obj := f()
	query := db.client.Collection(obj.Class()).Query
	return NewQueryFirestore(ctx, query, f)
}

//RunTransaction implement firestore run transaction
func (db *DBFirestore) RunTransaction(ctx context.Context, f func(tx protocol.Transaction) error) error {
	return db.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		tf := NewTransactionFirestore(ctx, db.client, tx)
		return f(tf)
	})
}
