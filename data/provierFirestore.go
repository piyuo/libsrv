package data

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

const credentialScopes = "https://www.googleapis.com/auth/datastore"

// ProviderFirestore implement google firestore
type ProviderFirestore struct {
	Provider
	//credentials *google.Credentials
	app *firebase.App
	ctx context.Context
}

// NewProviderFirestore provide new Provider for google firestore
func NewProviderFirestore() *ProviderFirestore {
	return &ProviderFirestore{}
}

//Initialize check env variable DATA_CRED to init google credentials for firestore
func (provider *ProviderFirestore) Initialize() {
	json := os.Getenv("SA")
	if json == "" {
		panic("Must set environment variable SA={google service account key} before using libsrv")
	}
	ctx := context.Background()
	creds, err := google.CredentialsFromJSON(ctx, []byte(json), credentialScopes)
	if err != nil {
		log.Printf("Failed to convert json to google credentials")
		panic(err)
	}
	serviceAccount := option.WithCredentials(creds)
	provider.app, err = firebase.NewApp(ctx, nil, serviceAccount)
	if err != nil {
		log.Printf("Failed to create firebase app")
		panic(err)
	}
	provider.ctx = ctx
}

//NewDB create db instance
func (provider *ProviderFirestore) NewDB() (IDB, error) {
	client, err := provider.app.Firestore(provider.ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create firestore client")
	}
	return NewDBFirestore(client), nil
}
