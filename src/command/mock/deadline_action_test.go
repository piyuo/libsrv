package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeadlineAction(t *testing.T) {
	assert := assert.New(t)
	action := &DeadlineAction{}
	response, err := action.Do(context.Background())
	assert.NotNil(err)
	assert.Nil(response)
}
