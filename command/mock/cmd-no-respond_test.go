package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdNoRespond(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	action := &CmdNoRespond{}
	response, err := action.Do(ctx)
	assert.Nil(err)
	assert.Nil(response)
}
