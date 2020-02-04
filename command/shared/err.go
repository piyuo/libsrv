package commands

import (
	"context"
)

// Main entry for client command execution
func (a *Err) Main(ctx context.Context) (interface{}, error) {
	panic("this is response, should not execute Main()")
}

//XXX_MapName override
func (a *Err) XXX_MapName() string {
	if a.Code == 0 {
		return "OK"
	}
	return "Err"
}
