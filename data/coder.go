package data

import "context"

// Coder is a collection of documents (shards) to realize code with high frequency.
//
type Coder interface {

	// CodeRX encode uint32 number into string, must used it in transaction with CodeWX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		code, err:= coder.CodeRX()
	//		So(err, ShouldBeNil)
	//		So(code, ShouldNotBeEmpty)
	//		err := coder.CodeWX()
	//	})
	//
	CodeRX(ctx context.Context) (string, error)

	// CodeWX commit CodeRX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		code, err:= coder.CodeRX()
	//		So(err, ShouldBeNil)
	//		So(code, ShouldNotBeEmpty)
	//		err := coder.CodeWX()
	//	})
	//
	CodeWX(ctx context.Context) error

	// Code16RX encode uint32 number into string, must used it in transaction with Code16WX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		code, err:= coder.Code16RX()
	//		So(err, ShouldBeNil)
	//		So(code, ShouldNotBeEmpty)
	//		err := coder.Code16WX()
	//	})
	//
	Code16RX(ctx context.Context) (string, error)

	// Code16WX commit Code16RX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		code, err:= coder.Code16RX()
	//		So(err, ShouldBeNil)
	//		So(code, ShouldNotBeEmpty)
	//		err := coder.Code16WX()
	//	})
	//
	Code16WX(ctx context.Context) error

	// Code64RX encode uint32 number into string, must used it in transaction with Code64WX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		code, err:= coder.Code64RX()
	//		So(err, ShouldBeNil)
	//		So(code, ShouldNotBeEmpty)
	//		err := coder.Code64WX()
	//	})
	//
	Code64RX(ctx context.Context) (string, error)

	// Code64WX commit with Code64RX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		code, err:= coder.Code64RX()
	//		So(err, ShouldBeNil)
	//		So(code, ShouldNotBeEmpty)
	//		err := coder.Code64WX()
	//	})
	//
	Code64WX(ctx context.Context) error

	// NumberRX prepare return unique but not serial number, must used it in transaction with NumberWX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		num, err:= coder.NumberRX()
	//		So(err, ShouldBeNil)
	//		So(num > 0, ShouldBeTrue)
	//		err := coder.NumberWX()
	//	})
	//
	NumberRX(ctx context.Context) (int64, error)

	// NumberWX commit NumberRX()
	//
	//	err = db.Transaction(ctx, func(ctx context.Context) error {
	//		num, err:= coder.NumberRX()
	//		So(err, ShouldBeNil)
	//		So(num > 0, ShouldBeTrue)
	//		err := coder.NumberWX()
	//	})
	//
	NumberWX(ctx context.Context) error

	// Clear all shards
	//
	//	err = coder.Clear(ctx)
	//
	Clear(ctx context.Context) error

	// ShardsCount returns shards count
	//
	//	count, err = coder.ShardsCount(ctx)
	//
	ShardsCount(ctx context.Context) (int, error)
}
