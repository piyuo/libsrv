package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFirestore(t *testing.T) {

	Convey("should create db ", t, func() {
		ctx := context.Background()
		db, err := firestoreNewDB(ctx)
		So(err, ShouldBeNil)
		So(db, ShouldNotBeNil)
	})
}
