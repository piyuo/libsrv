package data

import (
	"context"

	"cloud.google.com/go/firestore"
	gcp "github.com/piyuo/libsrv/secure/gcp"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// firestoreGlobalDB create global database instance
//
//	firestoreGlobalDB(ctx)
//
func firestoreGlobalDB(ctx context.Context) (DB, error) {
	cred, err := gcp.GlobalCredential(ctx)
	if err != nil {
		return nil, err
	}
	return firestoreNewDB(ctx, cred)
}

// firestoreRegionalDB create regional database instance
//
//	firestoreRegionalDB(ctx)
//
func firestoreRegionalDB(ctx context.Context) (DB, error) {
	cred, err := gcp.CurrentRegionalCredential(ctx)
	if err != nil {
		return nil, err
	}
	return firestoreNewDB(ctx, cred)
}

// firestoreNewDB create db instance
//
// firestoreNewDB(ctx)
//
func firestoreNewDB(ctx context.Context, cred *google.Credentials) (DB, error) {
	client, err := firestore.NewClient(ctx, cred.ProjectID, option.WithCredentials(cred))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create firestore client")
	}

	db := &DBFirestore{
		client: client,
	}
	return db, nil
}
