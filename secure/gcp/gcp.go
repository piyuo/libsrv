package gcp

import (
	"context"

	app "github.com/piyuo/go-libsrv/app"
	tools "github.com/piyuo/go-libsrv/tools"
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

	key := "log-gcp"
	scope := "https://www.googleapis.com/auth/cloud-platform"
	if logCredGlobal == nil {
		cred, err := createCredential(ctx, key, scope)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create  google log credential, "+key+".key")
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
	key := "data-gcp"
	scope := "https://www.googleapis.com/auth/datastore"
	if dataCredGlobal == nil {
		cred, err := createCredential(ctx, key, scope)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create google data  credential, "+key+".key")
		}
		dataCredGlobal = cred
	}
	return dataCredGlobal, nil
}

//createCredential base on key and scope
func createCredential(ctx context.Context, key, scope string) (*google.Credentials, error) {
	keyPath, err := app.KeyPath(key)
	if err != nil {
		return nil, errors.Wrap(err, key+".key not found")
	}
	jsonfile, err := tools.NewJSONFile(keyPath)
	if err != nil {
		return nil, errors.Wrap(err, "can no open key file "+keyPath)
	}
	defer jsonfile.Close()

	text, err := jsonfile.Text()
	if err != nil {
		return nil, errors.Wrap(err, " keyfile content maybe empty or wrong format. "+keyPath)
	}

	creds, err := google.CredentialsFromJSON(ctx, []byte(text), scope)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert json to google credentials.\n"+text)
	}
	return creds, nil
}
