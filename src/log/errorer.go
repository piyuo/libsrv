package log

import (
	"context"
)

// Errorer keep errors
//
type Errorer interface {

	// Close Errorer
	//
	//	defer errorer.Close()
	//
	Close()

	// Write error
	//
	//	errorer.write(ctx,here,"nil pointer",stack,"AAAAA")
	//
	Write(ctx context.Context, where, message, stack, errID string)
}
