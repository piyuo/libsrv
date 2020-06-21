package data

import (
	"context"
	"fmt"
	"testing"
	"time"

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

func TestReadWriteDelete(t *testing.T) {
	Convey("should read write file", t, func() {
		ctx := context.Background()
		storage, err := NewCloudstorage(ctx)
		bucketName := "mock-libsrv.piyuo.com"
		path := "text.txt"

		err = storage.AddBucket(ctx, bucketName, "US")
		So(err, ShouldBeNil)

		err = storage.WriteText(ctx, bucketName, path, "hi")
		So(err, ShouldBeNil)

		txt, err := storage.ReadText(ctx, bucketName, path)
		So(err, ShouldBeNil)
		So(txt, ShouldEqual, "hi")

		err = storage.Delete(ctx, bucketName, path)
		So(err, ShouldBeNil)

		err = storage.RemoveBucket(ctx, bucketName)
		So(err, ShouldBeNil)
	})
}

func TestCleanBucket(t *testing.T) {
	ctx := context.Background()
	storage, err := NewCloudstorage(ctx)
	bucketName := "mock-libsrv.piyuo.com"
	path := "text.txt"

	err = storage.AddBucket(ctx, bucketName, "US")
	//	So(err, ShouldBeNil)

	for i := 0; i < 100; i++ {
		err = storage.WriteText(ctx, bucketName, fmt.Sprintf("%v%v", path, i), fmt.Sprintf("hi %v", i))
		fmt.Printf("add object:%v\n", i)
	}

	err = storage.CleanBucket(ctx, bucketName, 25*time.Second)
	err = storage.RemoveBucket(ctx, bucketName)
	fmt.Printf("%v", err)
}
