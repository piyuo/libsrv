package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdRespond(t *testing.T) {
	t.Parallel()
 	assert := assert.New(t)
    ctx:=context.Background()
	action := &CmdRespond{}
    //  ctx = session.SetUserID(ctx, "user1")

    response, err := action.Do(ctx)
	assert.Nil(err)
	assert.NotNil(response)

    //  sr := response.(*StringResponse)
	//  assert.Equal("hi", sr.Text )
}
