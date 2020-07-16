package data

// SerialRef can generate serial Sequence in low frequency. 1 per second, use it with high frequency will cause error
//
type SerialRef interface {

	// NumberRX return sequence number, number is unique and serial, please be aware serial can only generate one sequence per second, use it with high frequency will cause error and  must used it in transaction with NumberWX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		num, err:= coder.NumberRX()
	//		So(err, ShouldBeNil)
	//		So(num, ShouldEqual,1)
	//		err := coder.NumberWX()
	//	})
	//
	NumberRX() (int64, error)

	// NumberWX commit NumberRX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		num, err:= coder.NumberRX()
	//		So(err, ShouldBeNil)
	//		So(num, ShouldEqual,1)
	//		err := coder.NumberWX()
	//	})
	//
	NumberWX() error
}
