package gcloud

import (
	"context"
	"testing"
	"time"

	"github.com/piyuo/libsrv/src/gaccount"
	"github.com/stretchr/testify/assert"
)

func TestGcloudCreateHTTPTask(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)

	schedule := time.Now().UTC().Add(30 * time.Second)
	err = CreateHTTPTask(ctx, cred, "piyuo-beta", "us-central1", "ci-queue", "http://notExist", schedule, []byte{})
	assert.Nil(err)

	err = CreateHTTPTask(ctx, cred, "notExist", "notExist", "notExist", "http://notExist", schedule, []byte{})
	assert.Contains(err.Error(), "not exist")
}
