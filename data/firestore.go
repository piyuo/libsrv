package data

import (
	"context"

	"cloud.google.com/go/firestore"
	log "github.com/piyuo/libsrv/log"
	gcp "github.com/piyuo/libsrv/secure/gcp"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

//firestoreNewDB create db instance
func firestoreNewDB(ctx context.Context) (DB, error) {
	cred, err := gcp.DataCredential(ctx)
	if err != nil {
		log.Alert(ctx, here, "failed to get firestore credential")
		return nil, err
	}

	client, err := firestore.NewClient(ctx, cred.ProjectID, option.WithCredentials(cred))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create firestore client")
	}

	db := &DBFirestore{
		client: client,
	}
	return db, nil
}
