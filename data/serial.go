package data

import (
	"context"
)

// SerialRef can generate serial number in low frequency. 1 per second, use it with high frequency will cause error
//
type SerialRef interface {

	// Number create unique serial number, please be aware serial can only generate one number per second, use it with high frequency will cause error
	//
	//	num, err := serial.Number(ctx, "sample-id")
	//	So(num, ShouldEqual, 1)
	//
	Number(ctx context.Context) (int64, error)
}
