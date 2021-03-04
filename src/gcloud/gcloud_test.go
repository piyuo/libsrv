package gcloud

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGcloudCreateHTTPTask(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	err := CreateHTTPTask(ctx, "http://it-is-not-exist.com", []byte{}, nil, 30*time.Second)
	assert.Nil(err)

	TestModeAlwaySuccess()
	err = CreateHTTPTask(ctx, "http://notExist", []byte{}, nil, 30*time.Second)
	assert.Nil(err)

	TestModeAlwayFail()
	err = CreateHTTPTask(ctx, "http://notExist", []byte{}, nil, 30*time.Second)
	assert.NotNil(err)

	TestModeBackNormal()

	//	err = CreateHTTPTask(ctx, "http://notExist", nil, nil)
	//	assert.Nil(err)
}
