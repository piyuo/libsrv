package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSerial(t *testing.T) {
	Convey("should query table", t, func() {
		ctx := context.Background()
		connG, err := NewSampleGlobalDB(ctx)
		defer connG.Close()
		So(err, ShouldBeNil)
		serialG := connG.Serial()
		So(serialG, ShouldNotBeNil)

		connR, err := NewSampleRegionalDB(ctx, "sample-database")
		defer connR.Close()
		So(err, ShouldBeNil)
		serialR := connR.Serial()
		So(serialR, ShouldNotBeNil)

		err = serialG.Delete(ctx, "sample")
		So(err, ShouldBeNil)
		err = serialR.Delete(ctx, "sample")
		So(err, ShouldBeNil)

		serialTest(ctx, serialG)
		serialTest(ctx, serialR)

		err = serialG.Delete(ctx, "sample")
		So(err, ShouldBeNil)
		err = serialR.Delete(ctx, "sample")
		So(err, ShouldBeNil)
	})
}

func serialTest(ctx context.Context, serial *SampleSerial) {

	id, err := serial.SampleID(ctx)
	So(err, ShouldBeNil)
	So(id, ShouldNotBeEmpty)

}
