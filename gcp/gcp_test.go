package gcp

import (
	"context"
	"testing"
	"time"

	"github.com/piyuo/libsrv/key"
	"github.com/piyuo/libsrv/region"
	"github.com/stretchr/testify/assert"
)

func TestCredential(t *testing.T) {
	assert := assert.New(t)

	//should create google credential
	bytes, err := key.BytesWithoutCache("gcloud.json")
	assert.Nil(err)

	cred, err := createCredential(context.Background(), bytes)
	assert.Nil(err)
	assert.NotNil(cred)

	//should keep global credential
	assert.Nil(globalCredential)
	cred, err = GlobalCredential(context.Background())
	assert.Nil(err)
	assert.NotNil(cred)
	assert.NotNil(globalCredential)
}

func TestDataCredentialByRegion(t *testing.T) {
	assert := assert.New(t)
	region.Current = "us"
	cred, err := RegionalCredential(context.Background())
	assert.Nil(err)
	assert.NotNil(cred)

	region.Current = "jp"
	cred, err = RegionalCredential(context.Background())
	assert.Nil(err)
	assert.NotNil(cred)

	region.Current = "be"
	cred, err = RegionalCredential(context.Background())
	assert.Nil(err)
	assert.NotNil(cred)
}

func TestCredentialWhenContextCanceled(t *testing.T) {
	assert := assert.New(t)
	dateline := time.Now().Add(time.Duration(1) * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), dateline)
	defer cancel()
	time.Sleep(time.Duration(2) * time.Millisecond)
	_, err := GlobalCredential(ctx)
	assert.NotNil(err)
	_, err = RegionalCredential(ctx)
	assert.NotNil(err)
}
