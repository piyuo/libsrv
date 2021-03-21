package db

import "context"

// Coder is a collection of documents (shards) to realize code with high frequency.
//
type Coder interface {

	// CodeRX encode uint32 number into string, must used it in transaction with CodeWX()
	//
	//	err := Transaction(ctx, func(ctx context.Context,tx Transaction) error {
	//		code, err:= coder.CodeRX(ctx,tx)
	//		err := coder.CodeWX(ctx,tx)
	//	})
	//
	CodeRX(ctx context.Context, transaction Transaction) (string, error)

	// CodeWX commit CodeRX()
	//
	//	err := Transaction(ctx, func(ctx context.Context,tx Transaction) error {
	//		code, err:= coder.CodeRX(ctx,tx)
	//		err := coder.CodeWX(ctx,tx)
	//	})
	//
	CodeWX(ctx context.Context, transaction Transaction) error

	// Code16RX encode uint16 number into string, must used it in transaction with CodeWX()
	//
	//	err := Transaction(ctx, func(ctx context.Context,tx Transaction) error {
	//		code, err:= coder.Code16RX(ctx,tx)
	//		err := coder.Code16WX(ctx,tx)
	//	})
	//
	Code16RX(ctx context.Context, transaction Transaction) (string, error)

	// Code16WX commit Code16RX()
	//
	//	err := Transaction(ctx, func(ctx context.Context,tx Transaction) error {
	//		code, err:= coder.Code16RX(ctx,tx)
	//		err := coder.Code16WX(ctx,tx)
	//	})
	//
	Code16WX(ctx context.Context, transaction Transaction) error

	// Code64RX encode uint32 number into string, must used it in transaction with Code64WX()
	//
	//	err := Transaction(ctx, func(ctx context.Context,tx Transaction) error {
	//		code, err:= coder.Code64RX(ctx,tx)
	//		err := coder.Code64WX(ctx,tx)
	//	})
	//
	Code64RX(ctx context.Context, transaction Transaction) (string, error)

	// Code64WX commit with Code64RX()
	//
	//	err := Transaction(ctx, func(ctx context.Context,tx Transaction) error {
	//		code, err:= coder.Code64RX(ctx,tx)
	//		err := coder.Code64WX(ctx,tx)
	//	})
	//
	Code64WX(ctx context.Context, transaction Transaction) error

	// NumberRX prepare return unique but not serial number, must used it in transaction with NumberWX()
	//
	//	err := Transaction(ctx, func(ctx context.Context,tx Transaction) error {
	//		num, err:= coder.NumberRX(ctx,tx)
	//		err := coder.NumberWX(ctx,tx)
	//	})
	//
	NumberRX(ctx context.Context, transaction Transaction) (int64, error)

	// NumberWX commit NumberRX()
	//
	//	err := Transaction(ctx, func(ctx context.Context,tx Transaction) error {
	//		num, err:= coder.NumberRX(ctx,tx)
	//		err := coder.NumberWX(ctx,tx)
	//	})
	//
	NumberWX(ctx context.Context, transaction Transaction) error

	// Delete delete coder
	//
	//	err = Delete(ctx)
	//
	Delete(ctx context.Context) error

	// ShardsCount returns shards count
	//
	//	count, err = coder.ShardsCount(ctx)
	//
	ShardsCount(ctx context.Context) (int, error)
}
