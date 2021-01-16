package log

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/identifier"
	"github.com/piyuo/libsrv/session"
	"github.com/stretchr/testify/assert"
)

func TestGCPErrorer(t *testing.T) {
	assert := assert.New(t)
	appName = "error-gcp_test"
	ctx := context.Background()
	ctx = session.SetUserID(ctx, "user1")
	errorer, err := NewGCPErrorer(ctx)
	assert.Nil(err)
	assert.NotNil(errorer)
	defer errorer.Close()
	stack := "at firstLine (a.js:3)\nat secondLine (b.js:3)"
	id := identifier.UUID()
	errorer.Write(ctx, "TestGCPLogger", "write error", stack, id)
}

func TestGCPEmptyStack(t *testing.T) {
	assert := assert.New(t)
	appName = "error-gcp_test"
	ctx := context.Background()
	errorer, err := NewGCPErrorer(ctx)
	assert.Nil(err)
	assert.NotNil(errorer)
	defer errorer.Close()
	id := identifier.UUID()
	errorer.Write(ctx, "TestGCPLogger", "write error", "", id)
}
