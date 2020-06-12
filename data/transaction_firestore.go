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
	ctx    context.Context
}

// NewTransactionFirestore is google firestore transaction
func NewTransactionFirestore(ctx context.Context, client *firestore.Client, tx *firestore.Transaction) *TransactionFirestore {
	return &TransactionFirestore{ctx: ctx, client: client, tx: tx}
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
