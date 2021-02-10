package gstorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/piyuo/libsrv/src/gaccount"
	"github.com/stretchr/testify/assert"
)

func TestNewGstorage(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	storage, err := New(ctx, cred)
	assert.Nil(err)
	assert.NotNil(storage)
}

func TestGstorageBucket(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	storage, err := New(ctx, cred)

	bucketName := "gstorage.piyuo.com"

	storage.DeleteBucket(ctx, bucketName)

	exist, err := storage.IsBucketExists(ctx, bucketName)
	assert.Nil(err)
	assert.False(exist)

	err = storage.CreateBucket(ctx, bucketName, "us-central1", "region")
	assert.Nil(err)

	exist, err = storage.IsBucketExists(ctx, bucketName)
	assert.Nil(err)
	assert.True(exist)

	err = storage.PublicBucket(ctx, bucketName)
	assert.Nil(err)

	err = storage.MakeBucketWebsite(ctx, bucketName, time.Hour*24*365, []string{"GET", "POST"}, []string{"*"}, []string{""})
	assert.Nil(err)

	err = storage.DeleteBucket(ctx, bucketName)
	assert.Nil(err)

	exist, err = storage.IsBucketExists(ctx, bucketName)
	assert.Nil(err)
	assert.False(exist)
}

func TestGstorageReadWriteDelete(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	storage, err := New(ctx, cred)
	bucketName := "gstorage.piyuo.com"
	filename := "a/b.txt"
	jsonname := "a/b.json"
	prefix := "a"

	err = storage.CreateBucket(ctx, bucketName, "us-central1", "region")
	assert.Nil(err)
	defer storage.DeleteBucket(ctx, bucketName)

	found, err := storage.IsFileExists(ctx, bucketName, filename)
	assert.Nil(err)
	assert.False(found)

	err = storage.WriteText(ctx, bucketName, filename, "hi")
	assert.Nil(err)

	found, err = storage.IsFileExists(ctx, bucketName, filename)
	assert.Nil(err)
	assert.True(found)

	txt, err := storage.ReadText(ctx, bucketName, filename)
	assert.Nil(err)
	assert.Equal("hi", txt)

	err = storage.DeleteFile(ctx, bucketName, filename)
	assert.Nil(err)

	// test json
	json := map[string]interface{}{
		"a": time.Now().UTC(),
	}
	err = storage.WriteJSON(ctx, bucketName, jsonname, json)
	assert.Nil(err)

	json, err = storage.ReadJSON(ctx, bucketName, jsonname)
	assert.Nil(err)
	assert.NotNil(json["a"])

	err = storage.DeleteFile(ctx, bucketName, jsonname)
	assert.Nil(err)

	// test delete files in dir
	err = storage.WriteText(ctx, bucketName, filename, "hi")
	assert.Nil(err)

	files, err := storage.ListFiles(ctx, bucketName, prefix, "")
	assert.Nil(err)
	assert.Equal(len(files), 1)

	err = storage.DeleteFiles(ctx, bucketName, prefix)
	assert.Nil(err)

	files, err = storage.ListFiles(ctx, bucketName, prefix, "")
	assert.Nil(err)
	assert.Equal(len(files), 0)

}

func TestGstorageCleanBucket(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	storage, err := New(ctx, cred)
	bucketName := "gstorage.piyuo.com"
	path := "TestCleanBucket.txt"

	err = storage.CreateBucket(ctx, bucketName, "us-central1", "region")
	assert.Nil(err)

	for i := 0; i < 1; i++ {
		err = storage.WriteText(ctx, bucketName, fmt.Sprintf("%v%v", path, i), fmt.Sprintf("hi %v", i))
		//fmt.Printf("add object:%v\n", i)
	}
	err = storage.CleanBucket(ctx, bucketName)
	assert.Nil(err)
	err = storage.DeleteBucket(ctx, bucketName)
	assert.Nil(err)
}
