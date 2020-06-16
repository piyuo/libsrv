package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewCloudstorage(t *testing.T) {
	Convey("should new cloudstorage", t, func() {
		storage, err := NewCloudstorage(context.Background())
		So(err, ShouldBeNil)
		So(storage, ShouldNotBeNil)
	})
}

func TestBucket(t *testing.T) {
	Convey("should new cloudstorage", t, func() {
		ctx := context.Background()
		storage, err := NewCloudstorage(ctx)

		bucketName := "mock-libsrv.piyuo.com"

		err = storage.RemoveBucket(ctx, bucketName)
		So(err, ShouldBeNil)

		exist, err := storage.IsBucketExist(ctx, bucketName)
		So(err, ShouldBeNil)
		So(exist, ShouldBeFalse)

		err = storage.AddBucket(ctx, bucketName, "US")
		So(err, ShouldBeNil)

		exist, err = storage.IsBucketExist(ctx, bucketName)
		So(err, ShouldBeNil)
		So(exist, ShouldBeTrue)

		err = storage.RemoveBucket(ctx, bucketName)
		So(err, ShouldBeNil)

		exist, err = storage.IsBucketExist(ctx, bucketName)
		So(err, ShouldBeNil)
		So(exist, ShouldBeFalse)
	})
}
