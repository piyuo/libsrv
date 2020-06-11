package gcp

import (
	"context"

	app "github.com/piyuo/libsrv/app"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
)

//save credential to global, so every cloud function can reuse this credential
var logCredGlobal *google.Credentials

//LogCredential provide google credential for log
func LogCredential(ctx context.Context) (*google.Credentials, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	key := "gcloud"
	scope := "https://www.googleapis.com/auth/cloud-platform"
	if logCredGlobal == nil {
		cred, err := createCredential(ctx, key, scope)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create google log credential, check /keys/"+key+".key exist")
		}
		logCredGlobal = cred
	}
	return logCredGlobal, nil
}

//save credential to global, so every cloud function can reuse this credential
var dataCredGlobal *google.Credentials

//DataCredential provide google credential for data
func DataCredential(ctx context.Context) (*google.Credentials, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	key := "gcloud"
	scope := "https://www.googleapis.com/auth/datastore"
	if dataCredGlobal == nil {
		cred, err := createCredential(ctx, key, scope)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create google data credential, check /keys/"+key+".key exist")
		}
		dataCredGlobal = cred
	}
	return dataCredGlobal, nil
}

//createCredential base on key and scope
//
//	cred, err := createCredential(context.Background(), "gcloud", "https://www.googleapis.com/auth/cloud-platform")
func createCredential(ctx context.Context, key, scope string) (*google.Credentials, error) {

	text, err := app.Key(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get key "+key)
	}

	creds, err := google.CredentialsFromJSON(ctx, []byte(text), scope)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert json to google credentials"+text)
	}
	return creds, nil
}
