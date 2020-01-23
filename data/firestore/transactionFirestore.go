package data

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// TransactionFirestore implement google firestore
type TransactionFirestore struct {
	Transaction
	client *firestore.Client
	tx     *firestore.Transaction
	ctx    context.Context
}

// NewTransactionFirestore is google firestore transaction
func NewTransactionFirestore(ctx context.Context, client *firestore.Client, tx *firestore.Transaction) *TransactionFirestore {
	return &TransactionFirestore{ctx: ctx, client: client, tx: tx}
}

//Get data object from data store, return ErrNotFound if object not exist
func (trans *TransactionFirestore) Get(obj IObject) error {
	id := obj.ID()
	if id == "" {
		return ErrNotFound
	}

	class := obj.Class()
	ref := trans.client.Collection(class).Doc(id)
	snapshot, err := trans.tx.Get(ref)
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

//Put data object into data store
func (trans *TransactionFirestore) Put(obj IObject) error {
	class := obj.Class()
	if obj.ID() == "" {
		ref := trans.client.Collection(class).NewDoc()
		obj.SetID(ref.ID)
	}
	ref := trans.client.Collection(class).Doc(obj.ID())
	err := trans.tx.Set(ref, obj)
	if err != nil {
		return errors.Wrap(err, "put object failed")
	}
	return nil
}

//Delete data object from firestore
func (trans *TransactionFirestore) Delete(obj IObject) error {
	id := obj.ID()
	class := obj.Class()
	ref := trans.client.Collection(class).Doc(id)
	err := trans.tx.Delete(ref)
	if err != nil {
		return errors.Wrap(err, "delete object failed")
	}
	return nil
}
