package data

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/util"
	"github.com/stretchr/testify/assert"
)

func TestDBInCanceledContext(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	ctxCanceled := util.CanceledCtx()

	dbR, _ := NewSampleRegionalDB(ctx)
	assert.NotNil(dbR.GetConnection())

	err := dbR.Transaction(ctxCanceled, func(ctx context.Context) error {
		return nil
	})
	assert.NotNil(err)
	err = dbR.BatchCommit(ctxCanceled)
	assert.NotNil(err)
}
