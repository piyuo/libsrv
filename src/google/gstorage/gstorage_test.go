package gstorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/piyuo/libsrv/src/google/gaccount"
	"github.com/stretchr/testify/assert"
)

var bucketName = "test-gstorage-1"

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

	err = storage.SetPageAndCORS(ctx, bucketName, "*")
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
}

func TestGstorageDeleteFiles(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	storage, err := New(ctx, cred)

	err = storage.CreateBucket(ctx, bucketName, "us-central1", "region")
	assert.Nil(err)
	defer storage.DeleteBucket(ctx, bucketName)

	err = storage.WriteText(ctx, bucketName, "a/b/c.txt", "hi")
	assert.Nil(err)
	err = storage.WriteText(ctx, bucketName, "a/b/c/d.txt", "hi")
	assert.Nil(err)
	err = storage.WriteText(ctx, bucketName, "c.txt", "hi")
	assert.Nil(err)

	err = storage.DeleteFiles(ctx, bucketName, "a/b")
	assert.Nil(err)

	found, err := storage.IsFileExists(ctx, bucketName, "a/b/c.txt")
	assert.Nil(err)
	assert.False(found)
	found, err = storage.IsFileExists(ctx, bucketName, "a/b/c/d.txt")
	assert.Nil(err)
	assert.False(found)

	found, err = storage.IsFileExists(ctx, bucketName, "c.txt")
	assert.Nil(err)
	assert.True(found)
}

func TestGstorageCleanBucket(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	storage, err := New(ctx, cred)
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

func TestGstorageSyncDir(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	storage, err := New(ctx, cred)

	err = storage.CreateBucket(ctx, bucketName, "us-central1", "region")
	assert.Nil(err)
	defer storage.DeleteBucket(ctx, bucketName)

	err = storage.WriteText(ctx, bucketName, "a/b/c.txt", "hi")
	assert.Nil(err)
	err = storage.WriteText(ctx, bucketName, "a/b/c/d.txt", "hi")
	assert.Nil(err)
	err = storage.WriteText(ctx, bucketName, "c.txt", "hi")
	assert.Nil(err)

	list := []string{}
	shouldDeleteFile := func(filename string) (bool, error) {
		list = append(list, filename)
		return true, nil
	}
	err = storage.Sync(ctx, bucketName, shouldDeleteFile)
	assert.Equal("a/b/c.txt", list[0])
	assert.Equal("a/b/c/d.txt", list[1])
	assert.Equal("c.txt", list[2])
}
