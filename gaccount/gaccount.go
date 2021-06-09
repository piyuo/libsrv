package gaccount

import (
	"context"
	"fmt"
	"sync"

	"github.com/piyuo/libsrv/env"
	"github.com/piyuo/libsrv/file"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
)

// Mock define key test flag
//
type Mock int8

const (
	// TestCredential use gcloud-test.json instead of gcloud.json
	//
	TestCredential Mock = iota
)

// forceTestCredential set to true will force using gcloud-test.json as credential
//
var forceTestCredential = false

// ForceTestCredential set to true will force using gcloud-test.json as credential
//
func ForceTestCredential(value bool) {
	forceTestCredential = value
}

// globalCredential keep global data credential to reuse in the future
//
var globalCredential *google.Credentials

// regionalCredentials keep regional credential to reuse in the future
//
var regionalCredentials map[string]*google.Credentials = make(map[string]*google.Credentials)

var regionalCredentialsMutex = sync.RWMutex{}

// googleMapKey keep google map api key
//
var googleMapKey string

// ClearCache clear credential cache
//
func ClearCache() {
	googleMapKey = ""
	globalCredential = nil
	regionalCredentials = make(map[string]*google.Credentials)
}

// GoogleMapKey provide google map key
//
//	key, err := GoogleMapKey(context.Background())
//
func GoogleMapKey(ctx context.Context) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	if googleMapKey == "" {
		var keyFile = "gmap.key"
		if forceTestCredential || ctx.Value(TestCredential) != nil {
			keyFile = "gmap-test.key"
		}
		text, err := file.KeyText(keyFile)
		if err != nil {
			return "", errors.Wrap(err, "get "+keyFile)
		}
		googleMapKey = text
	}
	return googleMapKey, nil
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
		ext := ""
		if forceTestCredential || ctx.Value(TestCredential) != nil {
			ext = "-test"
		}
		keyFilename := fmt.Sprintf("gcloud%s.json", ext)
		bytes, err := file.Key(keyFilename)
		if err != nil {
			return nil, errors.Wrap(err, "get key "+keyFilename)
		}
		cred, err := MakeCredential(ctx, bytes)
		if err != nil {
			return nil, errors.Wrapf(err, "make cred %v", keyFilename)
		}
		globalCredential = cred
	}
	return globalCredential, nil
}

// RegionalCredential provide google credential for regional database, region is set by os.Getenv("REGION")
//
//	cred, err := RegionalCredential(context.Background(), "us")
//
func RegionalCredential(ctx context.Context) (*google.Credentials, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	var cred = regionalCredentials[env.Region]
	if regionalCredentials[env.Region] == nil {
		ext := "-" + env.Region
		if forceTestCredential || ctx.Value(TestCredential) != nil {
			ext = "-test"
		}
		keyFilename := fmt.Sprintf("gcloud%s.json", ext)
		bytes, err := file.Key(keyFilename)
		if err != nil {
			return nil, errors.Wrap(err, "get key "+keyFilename)
		}
		cred, err := MakeCredential(ctx, bytes)
		if err != nil {
			return nil, errors.Wrapf(err, "make cred %v", keyFilename)
		}
		regionalCredentialsMutex.Lock()
		regionalCredentials[env.Region] = cred
		regionalCredentialsMutex.Unlock()
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

// CredentialFromFile get credential from key json, this function is used by ci
//
//	cred, err := CredentialFromFile(ctx,jsonFile)
//
func CredentialFromFile(ctx context.Context, jsonFile string) (*google.Credentials, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	bytes, err := file.ReadDirect(jsonFile)
	if err != nil {
		return nil, errors.Wrapf(err, "get json file %v", jsonFile)
	}
	cred, err := MakeCredential(ctx, bytes)
	if err != nil {
		return nil, err
	}
	return cred, nil
}

// AccountProjectFromFile return account and project from key file, this function is used by ci
//
//	account, project, err := AccountProjectFromFile(ctx,jsonFile)
//
func AccountProjectFromFile(ctx context.Context, jsonFile string) (string, string, error) {
	if ctx.Err() != nil {
		return "", "", ctx.Err()
	}

	j, err := file.ReadJSONDirect(jsonFile)
	if err != nil {
		return "", "", errors.Wrapf(err, "get json file %v", jsonFile)
	}
	return j["client_email"].(string), j["project_id"].(string), nil
}
