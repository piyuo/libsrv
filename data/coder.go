package data

import "context"

// CoderRef is a collection of documents (shards) to realize code with high frequency.
//
type CoderRef interface {

	// CodeRX encode uint32 number into string, must used it in transaction with CodeWX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		code, err:= coder.CodeRX()
	//		So(err, ShouldBeNil)
	//		So(code, ShouldNotBeEmpty)
	//		err := coder.CodeWX()
	//	})
	//
	CodeRX() (string, error)

	// CodeWX commit CodeRX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		code, err:= coder.CodeRX()
	//		So(err, ShouldBeNil)
	//		So(code, ShouldNotBeEmpty)
	//		err := coder.CodeWX()
	//	})
	//
	CodeWX() error

	// Code16RX encode uint32 number into string, must used it in transaction with Code16WX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		code, err:= coder.Code16RX()
	//		So(err, ShouldBeNil)
	//		So(code, ShouldNotBeEmpty)
	//		err := coder.Code16WX()
	//	})
	//
	Code16RX() (string, error)

	// Code16WX commit Code16RX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		code, err:= coder.Code16RX()
	//		So(err, ShouldBeNil)
	//		So(code, ShouldNotBeEmpty)
	//		err := coder.Code16WX()
	//	})
	//
	Code16WX() error

	// Code64RX encode uint32 number into string, must used it in transaction with Code64WX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		code, err:= coder.Code64RX()
	//		So(err, ShouldBeNil)
	//		So(code, ShouldNotBeEmpty)
	//		err := coder.Code64WX()
	//	})
	//
	Code64RX() (string, error)

	// Code64WX commit with Code64RX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		code, err:= coder.Code64RX()
	//		So(err, ShouldBeNil)
	//		So(code, ShouldNotBeEmpty)
	//		err := coder.Code64WX()
	//	})
	//
	Code64WX() error

	// NumberRX prepare return unique but not serial number, must used it in transaction with NumberWX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		num, err:= coder.NumberRX()
	//		So(err, ShouldBeNil)
	//		So(num > 0, ShouldBeTrue)
	//		err := coder.NumberWX()
	//	})
	//
	NumberRX() (int64, error)

	// NumberWX commit NumberRX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		num, err:= coder.NumberRX()
	//		So(err, ShouldBeNil)
	//		So(num > 0, ShouldBeTrue)
	//		err := coder.NumberWX()
	//	})
	//
	NumberWX() error

	// Reset reset code
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		err:= code.Reset(ctx)
	//	})
	//
	Reset(ctx context.Context) error
}
