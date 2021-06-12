package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdDeadline(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	action := &CmdDeadline{}
	response, err := action.Do(ctx)
	assert.NotNil(err)
	assert.Nil(response)
}
