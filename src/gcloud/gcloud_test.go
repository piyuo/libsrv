package gcloud

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGcloudCreateHTTPTask(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	err := CreateHTTPTask(ctx, "http://notExist", []byte{}, nil)
	assert.Nil(err)

	TestModeAlwaySuccess()
	err = CreateHTTPTask(ctx, "http://notExist", []byte{}, nil)
	assert.Nil(err)

	TestModeAlwayFail()
	err = CreateHTTPTask(ctx, "http://notExist", []byte{}, nil)
	assert.NotNil(err)

	TestModeBackNormal()
}
