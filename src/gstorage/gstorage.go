package gstorage

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"github.com/piyuo/libsrv/src/log"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const here = "gstorage"

// Gstorage is google cloud storage toolkit
//
//	before use this toolkit you need verify domain owner ship, please add service account in gclould.key to  https://www.google.com/webmasters/verification
//
type Gstorage interface {

	// AddBucket add cloud storage bucket
	//
	//	ctx := context.Background()
	//	storage, err := New(ctx)
	//	err = storage.AddBucket(ctx, "mock-libsrv.piyuo.com", "US")
	//
	AddBucket(ctx context.Context, bucketName, location string) error

	// RemoveBucket remove cloud storage bucket
	//
	//	ctx := context.Background()
	//	storage, err := New(ctx)
	//	err = storage.RemoveBucket(ctx, "mock-libsrv.piyuo.com")
	//
	RemoveBucket(ctx context.Context, bucketName string) error

	// CleanBucket remove all file in bucket
	//
	//	ctx := context.Background()
	//	storage, err := New(ctx)
	//	err = storage.RemoveBucket(ctx, "mock-libsrv.piyuo.com")
	//
	CleanBucket(ctx context.Context, bucketName string) error

	// IsBucketExists return true if bucket exist
	//
	//	bucketName := "mock-libsrv.piyuo.com"
	//	ctx := context.Background()
	//	storage, err := New(ctx)
	//	exist, err := storage.IsBucketExists(ctx, bucketName)
	//
	IsBucketExists(ctx context.Context, bucketName string) (bool, error)

	// IsFileExists return true if file exist
	//
	//	ctx := context.Background()
	//	storage, err := New(ctx)
	//	found,err = storage.IsFileExist(ctx, bucketName,"dirName", "fileName")
	//
	IsFileExists(ctx context.Context, bucketName, dirName, fileName string) (bool, error)

	// WriteText write text file to bucket
	//
	//	ctx := context.Background()
	//	storage, err := New(ctx)
	//	err = storage.AddBucket(ctx, bucketName, "US")
	//
	WriteText(ctx context.Context, bucketName, path, txt string) error

	// ReadText file from bucket
	//
	//	ctx := context.Background()
	//	storage, err := New(ctx)
	//	txt, err := storage.ReadText(ctx, bucketName, path)
	//	So(txt, ShouldEqual, "hi")
	//
	ReadText(ctx context.Context, bucketName, path string) (string, error)

	// DeleteFile file from bucket
	//
	//	ctx := context.Background()
	//	storage, err := New(ctx)
	//	err = storage.DeleteFile(ctx, bucketName, path)
	//
	DeleteFile(ctx context.Context, bucketName, path string) error

	//setCors
}

// Implementation is Cloudstorage implementation
//
type Implementation struct {
	Gstorage
	client    *storage.Client
	projectID string
}

// New create Cloudstorage
//
//	ctx := context.Background()
//	storage, err := New(ctx,cred)
//
func New(ctx context.Context, cred *google.Credentials) (Gstorage, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	client, err := storage.NewClient(ctx, option.WithCredentials(cred))
	if err != nil {
		return nil, err
	}

	cloudstorage := &Implementation{
		client:    client,
		projectID: cred.ProjectID,
	}
	return cloudstorage, nil
}

// AddBucket add cloud storage bucket
//
//	ctx := context.Background()
//	storage, err := New(ctx)
//	err = storage.AddBucket(ctx, "mock-libsrv.piyuo.com", "US")
//
func (impl *Implementation) AddBucket(ctx context.Context, bucketName, location string) error {

	exist, err := impl.IsBucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exist {
		bucket := impl.client.Bucket(bucketName)
		if err := bucket.Create(ctx, impl.projectID, &storage.BucketAttrs{
			Location: location,
		}); err != nil {
			return errors.Wrap(err, "failed to add bucket:"+bucketName)
		}

		log.Info(ctx, here, bucketName+" Bucket created")
	}
	return nil
}

// RemoveBucket remove cloud storage bucket
//
//	ctx := context.Background()
//	storage, err := New(ctx)
//	err = storage.RemoveBucket(ctx, "mock-libsrv.piyuo.com")
//
func (impl *Implementation) RemoveBucket(ctx context.Context, bucketName string) error {

	exist, err := impl.IsBucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if exist {
		bucket := impl.client.Bucket(bucketName)

		if err := impl.CleanBucket(ctx, bucketName); err != nil {
			return errors.Wrap(err, "failed to clean bucket before remove:"+bucketName)
		}

		if err := bucket.Delete(ctx); err != nil {
			return errors.Wrap(err, "failed to remove bucket:"+bucketName)
		}
		log.Info(ctx, here, bucketName+" Bucket deleted")
	}
	return nil
}

// IsBucketExists return true if bucket exist
//
//	bucketName := "mock-libsrv.piyuo.com"
//	ctx := context.Background()
//	storage, err := New(ctx)
//	exist, err := storage.IsBucketExist(ctx, bucketName)
//
func (impl *Implementation) IsBucketExists(ctx context.Context, bucketName string) (bool, error) {
	bucketIterator := impl.client.Buckets(ctx, impl.projectID)
	for {
		bucketAttrs, err := bucketIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, errors.Wrap(err, "failed iterator buckets")
		}

		if bucketAttrs.Name == bucketName {
			return true, nil
		}
	}
	return false, nil
}

// IsFileExists return true if file exist
//
//	ctx := context.Background()
//	storage, err := New(ctx)
//	found,err = storage.IsFileExist(ctx, bucketName, "dirName", "fileName")
//
func (impl *Implementation) IsFileExists(ctx context.Context, bucketName, dirName, fileName string) (bool, error) {
	bucket := impl.client.Bucket(bucketName)

	query := &storage.Query{Prefix: dirName}
	it := bucket.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, err
		}
		if attrs.Name == fileName {
			return true, nil
		}
	}
	return false, nil
}

// WriteText file to bucket
//
//	ctx := context.Background()
//	storage, err := New(ctx)
//	err = storage.AddBucket(ctx, bucketName, "US"
//
func (impl *Implementation) WriteText(ctx context.Context, bucketName, path, txt string) error {
	bucket := impl.client.Bucket(bucketName)

	wc := bucket.Object(path).NewWriter(ctx)
	_, err := io.WriteString(wc, txt)
	if err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}

// ReadText file from bucket
//
//	ctx := context.Background()
//	storage, err := New(ctx)
//	txt, err := storage.ReadText(ctx, bucketName, path
//	So(txt, ShouldEqual, "hi")
//
func (impl *Implementation) ReadText(ctx context.Context, bucketName, path string) (string, error) {
	bucket := impl.client.Bucket(bucketName)

	rc, err := bucket.Object(path).NewReader(ctx)
	if err != nil {
		return "", err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DeleteFile file from bucket
//
//	ctx := context.Background()
//	storage, err := New(ctx)
//	err = storage.DeleteFile(ctx, bucketName, path)
//
func (impl *Implementation) DeleteFile(ctx context.Context, bucketName, path string) error {

	o := impl.client.Bucket(bucketName).Object(path)
	if err := o.Delete(ctx); err != nil {
		return err
	}
	return nil
}

// CleanBucket remove all files in bucket, return true if still have file in bucket
// timeout in ms
//
//	ctx := context.Background()
//	storage, err := New(ctx)
//	err = storage.Delete(ctx, bucketName, path)
//
func (impl *Implementation) CleanBucket(ctx context.Context, bucketName string) error {
	bucket := impl.client.Bucket(bucketName)
	for {
		result, err := impl.RemoveObjects(ctx, bucket)
		if err != nil {
			return err
		}
		if result == true {
			return nil
		}
	}
}

// RemoveObjects remove objects max 1000, return true if object all deleted
//
func (impl *Implementation) RemoveObjects(ctx context.Context, bucket *storage.BucketHandle) (bool, error) {
	query := &storage.Query{}
	query.SetAttrSelection([]string{"Name"})

	i := 0
	it := bucket.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, err
		}
		fmt.Printf("delete object:%v\n", i)
		if err := bucket.Object(attrs.Name).Delete(ctx); err != nil {
			return false, err
		}
		i++
		if i >= 1000 {
			return false, nil
		}
	}
	return true, nil
}
