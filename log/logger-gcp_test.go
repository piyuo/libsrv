package log

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/env"
	"github.com/stretchr/testify/assert"
)

func TestGCPLogger(t *testing.T) {
	assert := assert.New(t)
	appName = "log-gcp_test"
	ctx := context.Background()
	ctx = env.SetUserID(ctx, "user1")
	logger, err := NewGCPLogger(ctx)
	assert.Nil(err)
	assert.NotNil(logger)
	defer logger.Close()

	logger.Write(ctx, DEBUG, here, "TestGCPLogger")

	shouldPrint = false
	//empty message will not log
	logger.Write(ctx, DEBUG, here, "no print console")

	//empty message will not log
	logger.Write(ctx, DEBUG, here, "")
}
