package data

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/src/util"
	"github.com/stretchr/testify/assert"
)

func TestDBInCanceledContext(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	ctxCanceled := util.CanceledCtx()

	g, err := NewSampleGlobalDB(ctx)
	assert.Nil(err)
	assert.NotNil(g.GetConnection())

	err = g.Transaction(ctxCanceled, func(ctx context.Context) error {
		return nil
	})
	assert.NotNil(err)
	err = g.BatchCommit(ctxCanceled)
	assert.NotNil(err)
}
