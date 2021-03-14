package gerror

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/src/env"
	"github.com/stretchr/testify/assert"
)

func TestGerror(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()

	client, err := NewClient(ctx)
	assert.Nil(err)
	assert.NotNil(client)
	defer close(ctx, client)

	ctx = env.SetUserID(ctx, "user1")
	env.AppName = "TestGerror"

	stack := "at firstLine (a.js:3)\nat secondLine (b.js:3)"
	Write(ctx, "hi", stack)

	// empty stack
	Write(ctx, "no stack", "")
}
