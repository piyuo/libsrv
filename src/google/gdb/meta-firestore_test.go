package gdb

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/src/identifier"
	"github.com/piyuo/libsrv/src/util"
	"github.com/stretchr/testify/assert"
)

func TestMeta(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	client := sampleClient()

	name := "test-meta-" + identifier.RandomString(8)
	shards := MetaFirestore{
		client:     client.(*ClientFirestore),
		id:         "id",
		collection: name,
		numShards:  0,
	}

	//check canceled ctx
	ctxCanceled := util.CanceledCtx()
	err := shards.check(ctxCanceled)
	assert.NotNil(err)

	//check empty id
	shards.id = ""
	err = shards.check(ctx)
	assert.NotNil(err)

	//check empty table name
	shards.id = "id"
	shards.collection = ""
	err = shards.check(ctx)
	assert.NotNil(err)
}
