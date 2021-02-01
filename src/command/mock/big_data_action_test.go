package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBigDataAction(t *testing.T) {
	assert := assert.New(t)
	action := &BigDataAction{}

	response, err := action.Do(context.Background())
	assert.Nil(err)
	assert.NotNil(response)
	//sr := response.(*StringResponse)
	//assert.Equal("hi", sr.Text )
}
