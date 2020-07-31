package gcp

import (
	"context"
	"os"

	"github.com/piyuo/libsrv/key"

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
		bytes, err := key.BytesWithoutCache("gcloud.json")
		if err != nil {
			return nil, errors.Wrap(err, "failed to get keys/gcloud.json")
		}
		cred, err := createCredential(ctx, bytes)
		if err != nil {
			return nil, err
		}
		globalCredential = cred
	}
	return globalCredential, nil
}

//regionalCredentials keep regional credential to reuse in the future
//
var regionalCredentials map[string]*google.Credentials = make(map[string]*google.Credentials)

// RegionalCredential provide google credential for regional database, region is set by os.Getenv("REGION")
//
//	cred, err := RegionalCredential(context.Background(), "us")
//
func RegionalCredential(ctx context.Context) (*google.Credentials, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	var cred = regionalCredentials[gcpRegion]
	if regionalCredentials[gcpRegion] == nil {
		bytes, err := key.BytesWithoutCache("/region/" + gcpRegion + ".json")
		if err != nil {
			return nil, errors.Wrap(err, "failed to get keys/region/"+gcpRegion+".json")
		}
		cred, err := createCredential(ctx, bytes)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create credential, check keys/region/"+gcpRegion+".json format is correct")
		}
		regionalCredentials[gcpRegion] = cred
		return cred, nil
	}
	return cred, nil
}

var gcpRegion = os.Getenv("REGION")

// createCredential create google credential from json bytes
//
//	cred, err := createCredential(context.Background(),bytes)
//
func createCredential(ctx context.Context, bytes []byte) (*google.Credentials, error) {

	creds, err := google.CredentialsFromJSON(ctx, bytes,
		"https://www.googleapis.com/auth/siteverification",        // log, error
		"https://www.googleapis.com/auth/cloud-platform",          // log, error
		"https://www.googleapis.com/auth/devstorage.full_control", // storage
		"https://www.googleapis.com/auth/datastore")
	if err != nil {
		return nil, err
	}
	return creds, nil
}
