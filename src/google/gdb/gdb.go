package gdb

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/piyuo/libsrv/src/db"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// NewClient create google db client
//
//	cred, err := gaccount.GlobalCredential(ctx)
//	return NewClient(ctx, cred)
//
func NewClient(ctx context.Context, cred *google.Credentials) (db.Client, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	firestoreClient, err := firestore.NewClient(ctx, cred.ProjectID, option.WithCredentials(cred))
	if err != nil {
		return nil, err
	}

	client := &ClientFirestore{
		firestoreClient: firestoreClient,
	}
	return client, nil
}
