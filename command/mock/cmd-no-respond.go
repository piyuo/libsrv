package mock

import (
	"context"
)

// Do comments
//
//	return OK if success
//	return "INVALID_XXX" when something wrong
//
func (c *CmdNoRespond) Do(ctx context.Context) (interface{}, error) {
	return nil, nil
}
