package data

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewDBWhenContextCanceled(t *testing.T) {
	Convey("should return error when context canceled", t, func() {
		dateline := time.Now().Add(time.Duration(1) * time.Millisecond)
		ctx, cancel := context.WithDeadline(context.Background(), dateline)
		defer cancel()
		time.Sleep(time.Duration(2) * time.Second)
		db, err := NewDB(ctx)
		So(err, ShouldNotBeNil)
		So(db, ShouldBeNil)
	})
}
