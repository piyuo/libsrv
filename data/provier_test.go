package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEnv(t *testing.T) {

	ctx := context.Background()
	db, _ := NewFirestoreDB(ctx)
	Convey("Fail to get db", t, func() {
		So(db, ShouldNotBeNil)
	})
}
