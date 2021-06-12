package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdResponse(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	action := &CmdResponse{}
	assert.NotNil(action)
}
