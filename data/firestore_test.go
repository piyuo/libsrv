package data

import (
	"context"
	"testing"

	gcp "github.com/piyuo/libsrv/secure/gcp"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFirestoreNewDB(t *testing.T) {
	Convey("should create db", t, func() {
		ctx := context.Background()
		cred, err := gcp.GlobalCredential(ctx)
		So(err, ShouldBeNil)
		db, err := firestoreNewDB(ctx, cred)
		So(err, ShouldBeNil)
		So(db, ShouldNotBeNil)
	})
}

func TestFirestoreGlobalDB(t *testing.T) {
	Convey("should create global db", t, func() {
		ctx := context.Background()
		db, err := firestoreGlobalDB(ctx)
		So(err, ShouldBeNil)
		So(db, ShouldNotBeNil)
	})
}

func TestFirestoreRegionalDB(t *testing.T) {
	Convey("should create regional db", t, func() {
		ctx := context.Background()
		db, err := firestoreRegionalDB(ctx)
		So(err, ShouldBeNil)
		So(db, ShouldNotBeNil)
	})
}
