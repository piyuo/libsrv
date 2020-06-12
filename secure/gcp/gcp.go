package gcp

import (
	"context"
	"os"

	app "github.com/piyuo/libsrv/app"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
)

// globalLogCredential keep global log credential to reuse in the future
//
var globalLogCredential *google.Credentials

// LogCredential provide google credential for log
//
//	cred, err := gcp.LogCredential(ctx)
//
func LogCredential(ctx context.Context) (*google.Credentials, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if globalLogCredential == nil {
		key, err := app.Key("gcloud")
		if err != nil {
			return nil, errors.Wrap(err, "failed to get keys/gcloud.key")
		}
		scope := "https://www.googleapis.com/auth/cloud-platform"
		cred, err := createCredential(ctx, key, scope)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create google log credential, check /keys/"+key+".key exist")
		}
		globalLogCredential = cred
	}
	return globalLogCredential, nil
}

//globalDataCredential keep global data credential to reuse in the future
//
var globalDataCredential *google.Credentials

// GlobalDataCredential provide google credential for global database
//
//	cred, err := GlobalDataCredential(context.Background())
//
func GlobalDataCredential(ctx context.Context) (*google.Credentials, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if globalDataCredential == nil {
		key, err := app.Key("gcloud")
		if err != nil {
			return nil, errors.Wrap(err, "failed to get keys/gcloud.key")
		}
		scope := "https://www.googleapis.com/auth/datastore"
		cred, err := createCredential(ctx, key, scope)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create credential, check /keys/"+key+".key format is correct")
		}
		globalDataCredential = cred
	}
	return globalDataCredential, nil
}

//globalDataCredential keep global data credential to reuse in the future
//
var regionalDataCredentials map[string]*google.Credentials = make(map[string]*google.Credentials)

// DataCredentialByRegion provide google credential for regional database
//
//	cred, err := DataCredentialByRegion(context.Background(), "us")
//
func DataCredentialByRegion(ctx context.Context, region string) (*google.Credentials, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	var cred = regionalDataCredentials[region]
	if regionalDataCredentials[region] == nil {
		key, err := app.RegionKey(region)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get keys/regions/"+region+".key")
		}
		scope := "https://www.googleapis.com/auth/datastore"
		cred, err := createCredential(ctx, key, scope)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create credential, check /keys/regions/"+key+".key format is correct")
		}
		regionalDataCredentials[region] = cred
		return cred, nil
	}
	return cred, nil
}

// RegionalDataCredential provide google credential for current region
//
func RegionalDataCredential(ctx context.Context) (*google.Credentials, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	region := os.Getenv("PIYUO_REGION")
	if region == "" {
		panic("need env PIYUO_APP=\"us\"")
	}
	return DataCredentialByRegion(ctx, region)
}

// createCredential base on key and scope
//
//	cred, err := createCredential(context.Background(), "gcloud", "https://www.googleapis.com/auth/cloud-platform")
//
func createCredential(ctx context.Context, key, scope string) (*google.Credentials, error) {

	creds, err := google.CredentialsFromJSON(ctx, []byte(key), scope)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert json to google credentials, "+key)
	}
	return creds, nil
}
