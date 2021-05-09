package gstorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/piyuo/libsrv/gaccount"
	"github.com/piyuo/libsrv/identifier"
	"github.com/stretchr/testify/assert"
)

func TestNewGstorage(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	storage, err := New(ctx, cred)
	assert.Nil(err)
	assert.NotNil(storage)
}

func TestGstorageBucket(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	storage, err := New(ctx, cred)
	assert.Nil(err)
	bucketName := "test-gstorage-bucket-" + identifier.RandomNumber(12)

	exist, err := storage.IsBucketExists(ctx, bucketName)
	assert.Nil(err)
	assert.False(exist)

	err = storage.CreateBucket(ctx, bucketName)
	assert.Nil(err)
	defer storage.DeleteBucket(ctx, bucketName)

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

	defer TestModeBackNormal()
	TestModeAlwaySuccess()

	err = storage.CreateBucket(ctx, bucketName)
	assert.Nil(err)
	exist, err = storage.IsBucketExists(ctx, bucketName)
	assert.Nil(err)
	assert.False(exist)

	err = storage.PublicBucket(ctx, bucketName)
	assert.Nil(err)

	err = storage.SetPageAndCORS(ctx, bucketName, "*")
	assert.Nil(err)

	err = storage.DeleteBucket(ctx, bucketName)
	assert.Nil(err)

	// delete not exists bucket will not error
	err = storage.DeleteBucket(ctx, bucketName)
	assert.Nil(err)

	TestModeAlwayFail()
	err = storage.CreateBucket(ctx, bucketName)
	assert.NotNil(err)

	exist, err = storage.IsBucketExists(ctx, bucketName)
	assert.False(exist)
	assert.NotNil(err)

	err = storage.PublicBucket(ctx, bucketName)
	assert.NotNil(err)

	err = storage.SetPageAndCORS(ctx, bucketName, "*")
	assert.NotNil(err)

	err = storage.DeleteBucket(ctx, bucketName)
	assert.NotNil(err)
}

func TestGstorageRW(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	storage, err := New(ctx, cred)
	assert.Nil(err)
	filename := "a/b.txt"
	jsonname := "a/b.json"
	prefix := "a"
	bucketName := "test-gstorage-rw-" + identifier.RandomNumber(12)

	err = storage.CreateBucket(ctx, bucketName)
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

	defer TestModeBackNormal()
	TestModeAlwaySuccess()
	_, err = storage.IsFileExists(ctx, bucketName, filename)
	assert.Nil(err)
	_, err = storage.ReadText(ctx, bucketName, filename)
	assert.Nil(err)
	_, err = storage.ReadJSON(ctx, bucketName, jsonname)
	assert.Nil(err)
	err = storage.WriteJSON(ctx, bucketName, jsonname, json)
	assert.Nil(err)
	err = storage.WriteText(ctx, bucketName, filename, "hi")
	assert.Nil(err)
	_, err = storage.ListFiles(ctx, bucketName, prefix, "")
	assert.Nil(err)
	err = storage.DeleteFile(ctx, bucketName, jsonname)
	assert.Nil(err)

	TestModeAlwayFail()
	_, err = storage.IsFileExists(ctx, bucketName, filename)
	assert.NotNil(err)
	_, err = storage.ReadText(ctx, bucketName, filename)
	assert.NotNil(err)
	_, err = storage.ReadJSON(ctx, bucketName, jsonname)
	assert.NotNil(err)
	err = storage.WriteJSON(ctx, bucketName, jsonname, json)
	assert.NotNil(err)
	err = storage.WriteText(ctx, bucketName, filename, "hi")
	assert.NotNil(err)
	_, err = storage.ListFiles(ctx, bucketName, prefix, "")
	assert.NotNil(err)
	err = storage.DeleteFile(ctx, bucketName, jsonname)
	assert.NotNil(err)
}

func TestGstorageDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	storage, err := New(ctx, cred)
	assert.Nil(err)
	bucketName := "test-gstorage-delete-" + identifier.RandomNumber(12)

	err = storage.CreateBucket(ctx, bucketName)
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

	defer TestModeBackNormal()
	TestModeAlwaySuccess()
	err = storage.DeleteFiles(ctx, bucketName, "a/b")
	assert.Nil(err)

	TestModeAlwayFail()
	err = storage.DeleteFiles(ctx, bucketName, "a/b")
	assert.NotNil(err)

}

func TestGstorageCleanBucket(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	storage, err := New(ctx, cred)
	assert.Nil(err)
	path := "TestCleanBucket.txt"
	bucketName := "test-gstorage-clean-" + identifier.RandomNumber(12)

	err = storage.CreateBucket(ctx, bucketName)
	assert.Nil(err)
	defer storage.DeleteBucket(ctx, bucketName)

	for i := 0; i < 1; i++ {
		err = storage.WriteText(ctx, bucketName, fmt.Sprintf("%v%v", path, i), fmt.Sprintf("hi %v", i))
		assert.Nil(err)
		//fmt.Printf("add object:%v\n", i)
	}
	err = storage.CleanBucket(ctx, bucketName)
	assert.Nil(err)

	defer TestModeBackNormal()
	TestModeAlwaySuccess()
	err = storage.CleanBucket(ctx, bucketName)
	assert.Nil(err)

	TestModeAlwayFail()
	err = storage.CleanBucket(ctx, bucketName)
	assert.NotNil(err)

}

func TestGstorageSyncDir(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	cred, err := gaccount.GlobalCredential(ctx)
	assert.Nil(err)
	storage, err := New(ctx, cred)
	assert.Nil(err)
	bucketName := "test-gstorage-sync-" + identifier.RandomNumber(12)

	err = storage.CreateBucket(ctx, bucketName)
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
	assert.Nil(err)
	assert.Equal("a/b/c.txt", list[0])
	assert.Equal("a/b/c/d.txt", list[1])
	assert.Equal("c.txt", list[2])

	defer TestModeBackNormal()
	TestModeAlwaySuccess()
	err = storage.Sync(ctx, bucketName, shouldDeleteFile)
	assert.Nil(err)

	TestModeAlwayFail()
	err = storage.Sync(ctx, bucketName, shouldDeleteFile)
	assert.NotNil(err)

}
