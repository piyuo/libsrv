package gdb

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/gaccount"
	"github.com/piyuo/libsrv/test"
	"github.com/stretchr/testify/assert"
)

func TestGdb(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	client, err := NewClient(ctx, cred)
	assert.Nil(err)
	assert.NotNil(client)
}

func TestInCanceledContext(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	ctxCanceled := test.CanceledContext()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	client, err := NewClient(ctxCanceled, cred)
	assert.NotNil(err)
	assert.Nil(client)
}

func TestNewClient(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client, err := NewClient(ctx, nil)
	assert.NotNil(err)
	assert.Nil(client)
}
