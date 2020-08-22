package data

import (
	"context"
	"testing"

	util "github.com/piyuo/libsrv/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDBNamespace(t *testing.T) {
	Convey("Should create and delete namespace", t, func() {
		ctx := context.Background()

		dbR, err := NewSampleRegionalDB(ctx, "sample-namespace")
		So(err, ShouldBeNil)
		So(dbR.GetConnection(), ShouldNotBeNil)

		err = dbR.CreateNamespace(ctx)
		So(err, ShouldBeNil)

		exist, err := dbR.IsNamespaceExist(ctx)
		So(err, ShouldBeNil)
		So(exist, ShouldBeTrue)

		err = dbR.ClearNamespace(ctx)
		So(err, ShouldBeNil)

		exist, err = dbR.IsNamespaceExist(ctx)
		So(err, ShouldBeNil)
		So(exist, ShouldBeFalse)

		err = dbR.CreateNamespace(ctx)
		So(err, ShouldBeNil)

		exist, err = dbR.IsNamespaceExist(ctx)
		So(err, ShouldBeNil)
		So(exist, ShouldBeTrue)

		err = dbR.DeleteNamespace(ctx)
		So(err, ShouldBeNil)

		exist, err = dbR.IsNamespaceExist(ctx)
		So(err, ShouldBeNil)
		So(exist, ShouldBeFalse)
	})
}

func TestDBInCanceledContext(t *testing.T) {
	Convey("Should return error if ctx if canceled", t, func() {
		ctx := context.Background()
		ctxCanceled := util.CanceledCtx()

		dbR, _ := NewSampleRegionalDB(ctx, "sample-namespace")
		So(dbR.GetConnection(), ShouldNotBeNil)

		exist, err := dbR.IsNamespaceExist(ctxCanceled)
		So(err, ShouldNotBeNil)
		So(exist, ShouldBeFalse)

		err = dbR.ClearNamespace(ctxCanceled)
		So(err, ShouldNotBeNil)

		err = dbR.DeleteNamespace(ctxCanceled)
		So(err, ShouldNotBeNil)
		err = dbR.CreateNamespace(ctxCanceled)
		So(err, ShouldNotBeNil)
		err = dbR.Transaction(ctxCanceled, func(ctx context.Context) error {
			return nil
		})
		So(err, ShouldNotBeNil)
		err = dbR.BatchCommit(ctxCanceled)
		So(err, ShouldNotBeNil)
	})
}
