package gaccount

import (
	"context"

	"github.com/piyuo/libsrv/src/file"
	"github.com/piyuo/libsrv/src/key"
	"github.com/piyuo/libsrv/src/region"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
)

// testMode set to true will force gaccount use gcloud-test.json as key
//
var testMode = false

// TestMode set to true will let every function run success
//
func TestMode(enabled bool) {
	testMode = enabled
}

// ClearCache clear credential cache
//
func ClearCache() {
	globalCredential = nil
	regionalCredentials = make(map[string]*google.Credentials)
}

//globalCredential keep global data credential to reuse in the future
//
var globalCredential *google.Credentials

// CreateCredential create credential from key
//
//	cred, err := CreateCredential(ctx,"master/gcloud.json")
//
func CreateCredential(ctx context.Context, keyName string) (*google.Credentials, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	bytes, err := key.BytesWithoutCache(keyName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get keys/"+keyName)
	}
	cred, err := MakeCredential(ctx, bytes)
	if err != nil {
		return nil, err
	}
	return cred, nil
}

// GlobalCredential provide google credential for project
//
//	cred, err := GlobalCredential(context.Background())
//
func GlobalCredential(ctx context.Context) (*google.Credentials, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if globalCredential == nil {
		keyFile := "gcloud.json"
		if testMode {
			keyFile = "gcloud-test.json"
		}

		bytes, err := key.BytesWithoutCache(keyFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get keys/"+keyFile)
		}
		cred, err := MakeCredential(ctx, bytes)
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
	var cred = regionalCredentials[region.Current]
	if regionalCredentials[region.Current] == nil {
		keyFile := "gcloud-" + region.Current + ".json"
		if testMode {
			keyFile = "gcloud-test.json"
		}

		bytes, err := key.BytesWithoutCache(keyFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get keys/"+keyFile)
		}
		cred, err := MakeCredential(ctx, bytes)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create credential, check format on:"+keyFile)
		}
		regionalCredentials[region.Current] = cred
		return cred, nil
	}
	return cred, nil
}

// MakeCredential create google credential from json bytes
//
//	cred, err := MakeCredential(context.Background(),bytes)
//
func MakeCredential(ctx context.Context, bytes []byte) (*google.Credentials, error) {
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

// CredentialFromFile get credential from key json
//
//	cred, err := CredentialFromFile(ctx,jsonFile)
//
func CredentialFromFile(ctx context.Context, jsonFile string) (*google.Credentials, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	bytes, err := file.Read(jsonFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get json file:"+jsonFile)
	}
	cred, err := MakeCredential(ctx, bytes)
	if err != nil {
		return nil, err
	}
	return cred, nil
}
