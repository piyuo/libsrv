package gaccount

import (
	"context"
	"testing"
	"time"

	"github.com/piyuo/libsrv/env"
	"github.com/piyuo/libsrv/file"
	"github.com/stretchr/testify/assert"
)

func TestCredentials(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	//should create google credential
	bytes, err := file.Key("gcloud.json")
	assert.Nil(err)
	ctx := context.Background()
	cred, err := MakeCredential(ctx, bytes)
	assert.Nil(err)
	assert.NotNil(cred)

	//should keep global credential
	cred, err = GlobalCredential(ctx)
	assert.Nil(err)
	assert.NotNil(cred)
	assert.NotNil(globalCredential)
}

func TestGoogleMapApiKey(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	key, err := GoogleMapApiKey(ctx)
	assert.Nil(err)
	assert.NotEmpty(key)
	// from cache
	key2, err := GoogleMapApiKey(ctx)
	assert.Nil(err)
	assert.Equal(key, key2)
}

func TestDataCredentialByRegion(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	bak := env.Region
	env.Region = "us"
	ctx := context.Background()
	cred, err := RegionalCredential(ctx)
	assert.Nil(err)
	assert.NotNil(cred)

	env.Region = "jp"
	cred, err = RegionalCredential(ctx)
	assert.Nil(err)
	assert.NotNil(cred)

	env.Region = "be"
	cred, err = RegionalCredential(ctx)
	assert.Nil(err)
	assert.NotNil(cred)
	env.Region = bak
}

func TestCredentialWhenContextCanceled(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	deadline := time.Now().Add(time.Duration(1) * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	time.Sleep(time.Duration(2) * time.Millisecond)
	_, err := GlobalCredential(ctx)
	assert.NotNil(err)
	_, err = RegionalCredential(ctx)
	assert.NotNil(err)
}

func TestContextTestCredential(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ClearCache()
	ctx := context.WithValue(context.Background(), TestCredential, "")

	cred, err := GlobalCredential(ctx)
	assert.Nil(err)
	cred2, err := RegionalCredential(ctx)
	assert.Nil(err)
	assert.Equal(cred.ProjectID, cred2.ProjectID)
}

func TestForceTestCredential(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ClearCache()
	ctx := context.Background()
	ForceTestCredential(true)
	defer ForceTestCredential(true)

	cred, err := GlobalCredential(ctx)
	assert.Nil(err)
	cred2, err := RegionalCredential(ctx)
	assert.Nil(err)
	assert.Equal(cred.ProjectID, cred2.ProjectID)
}

func TestCredentialFromFile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	assert := assert.New(t)
	keyPath := "../../keys/gcloud-test.json"
	cred, err := CredentialFromFile(ctx, keyPath)
	assert.Nil(err)
	assert.NotNil(cred)

	cred, err = CredentialFromFile(ctx, "notExists")
	assert.NotNil(err)
	assert.Nil(cred)
}

func TestAccountProjectFromFile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	assert := assert.New(t)
	keyPath := "../../keys/gcloud-test.json"
	account, project, err := AccountProjectFromFile(ctx, keyPath)
	assert.Nil(err)
	assert.NotEmpty(account)
	assert.NotEmpty(project)
}
