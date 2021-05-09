package gaccount

import (
	"context"
	"testing"
	"time"

	"github.com/piyuo/libsrv/env"
	"github.com/piyuo/libsrv/key"
	"github.com/stretchr/testify/assert"
)

func TestGaccountCredential(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	//should create google credential
	bytes, err := key.BytesWithoutCache("gcloud.json")
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

func TestGaccountCreateCredential(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := CreateCredential(ctx, "gcloud.json")
	assert.Nil(err)
	assert.NotNil(cred)
	cred, err = CreateCredential(ctx, "notExist.json")
	assert.NotNil(err)
	assert.Nil(cred)
}

func TestGaccountDataCredentialByRegion(t *testing.T) {
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

func TestGaccountCredentialWhenContextCanceled(t *testing.T) {
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

func TestGaccountTestMode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	assert := assert.New(t)
	ClearCache()
	UseTestCredential(true)
	defer UseTestCredential(false)

	cred, err := GlobalCredential(ctx)
	assert.Nil(err)

	cred2, err := RegionalCredential(ctx)
	assert.Nil(err)

	assert.Equal(cred.ProjectID, cred2.ProjectID)
}

func TestGaccountCredentialFromFile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	assert := assert.New(t)
	cred, err := CredentialFromFile(ctx, "notExists")
	assert.NotNil(err)
	assert.Nil(cred)
}
