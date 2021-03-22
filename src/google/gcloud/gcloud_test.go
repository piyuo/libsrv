package gcloud

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGcloudCreateHTTPTask(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()

	err := CreateHTTPTask(ctx, "task", "http://it-is-not-exist.com", []byte{}, nil)
	assert.Nil(err)

	TestModeAlwaySuccess()
	err = CreateHTTPTask(ctx, "task", "http://notExist", []byte{}, nil)
	assert.Nil(err)

	TestModeAlwayFail()
	err = CreateHTTPTask(ctx, "task", "http://notExist", []byte{}, nil)
	assert.NotNil(err)

	TestModeBackNormal()
}
