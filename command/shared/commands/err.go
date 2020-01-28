package commands

import (
	"context"
)

// Main entry for client command execution
func (a *Err) Main(ctx context.Context) (interface{}, error) {
	panic("this is response, should not execute Main()")
}
