package data

import (
	"context"
	"testing"

	util "github.com/piyuo/libsrv/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDBInCanceledContext(t *testing.T) {
	Convey("Should return error if ctx if canceled", t, func() {
		ctx := context.Background()
		ctxCanceled := util.CanceledCtx()

		dbR, _ := NewSampleRegionalDB(ctx, "sample-namespace")
		err := dbR.DeleteNamespace(ctxCanceled)
		So(err, ShouldNotBeNil)
		err = dbR.CreateNamespace(ctxCanceled)
		So(err, ShouldNotBeNil)
		err = dbR.Transaction(ctxCanceled, func(ctx context.Context) error {
			return nil
		})
		So(err, ShouldNotBeNil)
	})
}
