package data

import (
	"context"
	"testing"

	"github.com/piyuo/libsrv/util"
	"github.com/stretchr/testify/assert"
)

func TestShardsFirestore(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	dbG, dbR := createSampleDB()
	defer removeSampleDB(dbG, dbR)

	shards := MetaFirestore{
		conn:      dbG.Connection.(*ConnectionFirestore),
		id:        "id",
		tableName: "tablename",
		numShards: 0,
	}
	id := shards.errorID()
	assert.Equal("tablename-id", id)

	//check canceled ctx
	ctxCanceled := util.CanceledCtx()
	err := shards.assert(ctxCanceled)
	assert.NotNil(err)

	//check empty id
	shards.id = ""
	err = shards.assert(ctx)
	assert.NotNil(err)

	//check empty table name
	shards.id = "id"
	shards.tableName = ""
	err = shards.assert(ctx)
	assert.NotNil(err)
}
