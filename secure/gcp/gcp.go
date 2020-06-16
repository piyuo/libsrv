package gcp

import (
	"context"
	"os"

	app "github.com/piyuo/libsrv/app"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
)

//globalCredential keep global data credential to reuse in the future
//
var globalCredential *google.Credentials

// GlobalCredential provide google credential for project
//
//	cred, err := GlobalCredential(context.Background())
//
func GlobalCredential(ctx context.Context) (*google.Credentials, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if globalCredential == nil {
		key, err := app.Key("gcloud")
		if err != nil {
			return nil, errors.Wrap(err, "failed to get keys/gcloud.key")
		}
		cred, err := createCredential(ctx, key,
			"https://www.googleapis.com/auth/cloud-platform",          // log, error
			"https://www.googleapis.com/auth/devstorage.full_control", // storage
			"https://www.googleapis.com/auth/datastore")               // firestore
		if err != nil {
			return nil, errors.Wrap(err, "failed to create credential, check /keys/"+key+".key format is correct")
		}
		globalCredential = cred
	}
	return globalCredential, nil
}

//regionalCredentials keep regional credential to reuse in the future
//
var regionalCredentials map[string]*google.Credentials = make(map[string]*google.Credentials)

// RegionalCredential provide google credential for regional database
//
//	cred, err := RegionalCredential(context.Background(), "us")
//
func RegionalCredential(ctx context.Context, region string) (*google.Credentials, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	var cred = regionalCredentials[region]
	if regionalCredentials[region] == nil {
		key, err := app.RegionKey(region)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get keys/regions/"+region+".key")
		}
		scope := "https://www.googleapis.com/auth/datastore"
		cred, err := createCredential(ctx, key, scope)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create credential, check /keys/regions/"+key+".key format is correct")
		}
		regionalCredentials[region] = cred
		return cred, nil
	}
	return cred, nil
}

// CurrentRegionalCredential provide google credential for current region
//
func CurrentRegionalCredential(ctx context.Context) (*google.Credentials, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	region := os.Getenv("PIYUO_REGION")
	if region == "" {
		panic("need env PIYUO_APP=\"us\"")
	}
	return RegionalCredential(ctx, region)
}

// createCredential base on key and scope
//
//	cred, err := createCredential(context.Background(), "gcloud", "https://www.googleapis.com/auth/cloud-platform")
//
func createCredential(ctx context.Context, key string, scope ...string) (*google.Credentials, error) {

	creds, err := google.CredentialsFromJSON(ctx, []byte(key), scope...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert json to google credentials, "+key)
	}
	return creds, nil
}
