package data

import "context"

// SerialRef can generate serial Sequence in low frequency. 1 per second, use it with high frequency will cause error
//
type SerialRef interface {

	// NumberRX return sequence number, number is unique and serial, start from 1 to 9,223,372,036,854,775,807, please be aware serial can only generate one sequence per second, use it with high frequency will cause error and  must used it in transaction with NumberWX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		num, err:= serial.NumberRX()
	//		So(err, ShouldBeNil)
	//		So(num, ShouldEqual,1)
	//		err := serial.NumberWX()
	//	})
	//
	NumberRX() (int64, error)

	// NumberWX commit NumberRX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		num, err:= serial.NumberRX()
	//		So(err, ShouldBeNil)
	//		So(num, ShouldEqual,1)
	//		err := serial.NumberWX()
	//	})
	//
	NumberWX() error

	// Reset reset sequence number
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		err:= serial.Reset(ctx)
	//	})
	//
	Reset(ctx context.Context) error
}
