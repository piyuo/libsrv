package gdb

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/src/google/gaccount"
	"github.com/piyuo/libsrv/src/util"
	"github.com/stretchr/testify/assert"
)

func TestGdb(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.NotNil(err)

	client, err := NewClient(ctx, cred)
	assert.Nil(err)
	assert.NotNil(client)
}

func TestGdbInCanceledContext(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	ctxCanceled := util.CanceledCtx()

	cred, err := gaccount.GlobalCredential(ctx)
	assert.NotNil(err)

	client, err := NewClient(ctxCanceled, cred)
	assert.NotNil(err)
	assert.Nil(client)
}
