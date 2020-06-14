package data

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

// TransactionFirestore implement google firestore
type TransactionFirestore struct {
	AbstractTransaction
	client *firestore.Client
	tx     *firestore.Transaction
}

// NewTransactionFirestore is google firestore transaction
func NewTransactionFirestore(ctx context.Context, client *firestore.Client, tx *firestore.Transaction) *TransactionFirestore {
	return &TransactionFirestore{client: client, tx: tx}
}

//Get data object from data store, return ErrNotFound if object not exist
func (trans *TransactionFirestore) Get(ctx context.Context, obj Object) error {
	id := obj.ID()
	if id == "" {
		return errors.New("get object need object  have ID")
	}

	modelName := obj.ModelName()
	ref := trans.client.Collection(modelName).Doc(id)
	snapshot, err := trans.tx.Get(ref)

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

//Put data object into data store
func (trans *TransactionFirestore) Put(ctx context.Context, obj Object) error {
	modelName := obj.ModelName()
	if obj.ID() == "" {
		ref := trans.client.Collection(modelName).NewDoc()
		obj.SetID(ref.ID)
	}
	ref := trans.client.Collection(modelName).Doc(obj.ID())
	err := trans.tx.Set(ref, obj)
	if err != nil {
		return errors.Wrap(err, "put object failed")
	}
	return nil
}

//Delete data object from firestore
func (trans *TransactionFirestore) Delete(ctx context.Context, obj Object) error {
	id := obj.ID()
	modelName := obj.ModelName()
	ref := trans.client.Collection(modelName).Doc(id)
	err := trans.tx.Delete(ref)
	if err != nil {
		return errors.Wrap(err, "delete object failed")
	}
	return nil
}

// ShortID create unique serial number, please be aware serial can only generate one number per second and use with transation to ensure unique
//
//	id,err = db.ShortID(ctx, "myID")
//
func (trans *TransactionFirestore) ShortID(ctx context.Context, name string) (*ShortID, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	shortID := &ShortID{
		ID: 1,
	}
	docRef := trans.client.Collection("shortID").Doc(name)
	snapshot, err := trans.tx.Get(docRef)

	if snapshot != nil && !snapshot.Exists() {
		err := trans.tx.Set(docRef, shortID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to put id:"+name)
		}
		return shortID, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get id:"+name)
	}

	if err := snapshot.DataTo(shortID); err != nil {
		return nil, errors.Wrap(err, "failed convert to id:"+name)
	}
	shortID.ID++

	trans.tx.Update(docRef, []firestore.Update{
		{Path: "Next", Value: firestore.Increment(1)},
	})

	return shortID, nil
}
