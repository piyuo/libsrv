package cloudstorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCloudstorage(t *testing.T) {
	assert := assert.New(t)
	storage, err := NewCloudstorage(context.Background())
	assert.Nil(err)
	assert.NotNil(storage)
}

func TestBucket(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	storage, err := NewCloudstorage(ctx)

	bucketName := "mock-libsrv.piyuo.com"

	err = storage.RemoveBucket(ctx, bucketName)
	assert.Nil(err)

	exist, err := storage.IsBucketExist(ctx, bucketName)
	assert.Nil(err)
	assert.False(exist)

	err = storage.AddBucket(ctx, bucketName, "US")
	assert.Nil(err)

	exist, err = storage.IsBucketExist(ctx, bucketName)
	assert.Nil(err)
	assert.True(exist)

	err = storage.RemoveBucket(ctx, bucketName)
	assert.Nil(err)

	exist, err = storage.IsBucketExist(ctx, bucketName)
	assert.Nil(err)
	assert.False(exist)
}

func TestReadWriteDelete(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	storage, err := NewCloudstorage(ctx)
	bucketName := "mock-libsrv.piyuo.com"
	path := "TestReadWriteDelete.txt"

	err = storage.AddBucket(ctx, bucketName, "US")
	assert.Nil(err)

	err = storage.WriteText(ctx, bucketName, path, "hi")
	assert.Nil(err)

	txt, err := storage.ReadText(ctx, bucketName, path)
	assert.Nil(err)
	assert.Equal("hi", txt)

	err = storage.Delete(ctx, bucketName, path)
	assert.Nil(err)

	err = storage.RemoveBucket(ctx, bucketName)
	assert.Nil(err)
}

func TestCleanBucket(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	storage, err := NewCloudstorage(ctx)
	bucketName := "mock-libsrv.piyuo.com"
	path := "TestCleanBucket.txt"

	err = storage.AddBucket(ctx, bucketName, "US")
	assert.Nil(err)

	for i := 0; i < 1; i++ {
		err = storage.WriteText(ctx, bucketName, fmt.Sprintf("%v%v", path, i), fmt.Sprintf("hi %v", i))
		//fmt.Printf("add object:%v\n", i)
	}
	err = storage.CleanBucket(ctx, bucketName, 25*time.Second)
	assert.Nil(err)
	err = storage.RemoveBucket(ctx, bucketName)
	assert.Nil(err)
}
