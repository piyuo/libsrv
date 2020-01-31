package data

import (
	"context"

	firebase "firebase.google.com/go"
	log "github.com/piyuo/go-libsrv/log"
	gcp "github.com/piyuo/go-libsrv/secure/gcp"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

var firebaseApp *firebase.App

//NewFirestoreDB create db instance
func NewFirestoreDB(ctx context.Context) (DB, error) {
	if firebaseApp == nil {
		cred, err := gcp.DataCredential(ctx)
		if err != nil {
			log.Alert(ctx, here, "database operation failed to get data google credential")
			return nil, err
		}

		firebaseApp, err = firebase.NewApp(ctx, nil, option.WithCredentials(cred))
		if err != nil {
			log.Alert(ctx, here, "database operation failed to create firebase app")
			return nil, err
		}
	}

	client, err := firebaseApp.Firestore(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create firestore client")
	}

	db := &DBFirestore{
		client: client,
	}
	return db, nil
}
