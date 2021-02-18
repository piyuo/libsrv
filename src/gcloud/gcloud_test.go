package gcloud

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/src/gaccount"
	"github.com/stretchr/testify/assert"
)

func TestGcloudCreateHTTPTask(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	err = CreateHTTPTask(ctx, cred, "notExist", "notExist", "notExist", "http://notExist", []byte{})
	assert.Contains(err.Error(), "not exist")
}
