package cloudstorage

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
	"github.com/piyuo/libsrv/gcp"
	"github.com/piyuo/libsrv/log"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const here = "cloudstorage"

// Cloudstorage is google cloud storage toolkit
//
//	before use this toolkit you need verify domain owner ship, please add service account in gclould.key to  https://www.google.com/webmasters/verification
//
type Cloudstorage interface {

	// AddBucket add cloud storage bucket
	//
	//	ctx := context.Background()
	//	storage, err := NewCloudstorage(ctx)
	//	err = storage.AddBucket(ctx, "mock-libsrv.piyuo.com", "US")
	//
	AddBucket(ctx context.Context, bucketName, location string) error

	// RemoveBucket remove cloud storage bucket
	//
	//	ctx := context.Background()
	//	storage, err := NewCloudstorage(ctx)
	//	err = storage.RemoveBucket(ctx, "mock-libsrv.piyuo.com")
	//
	RemoveBucket(ctx context.Context, bucketName string) error

	// CleanBucket remove all file in bucket
	//
	//	ctx := context.Background()
	//	storage, err := NewCloudstorage(ctx)
	//	err = storage.RemoveBucket(ctx, "mock-libsrv.piyuo.com")
	//
	CleanBucket(ctx context.Context, bucketName string, timeout time.Duration) error

	// IsBucketExist return true if bucket exist
	//
	//	bucketName := "mock-libsrv.piyuo.com"
	//	ctx := context.Background()
	//	storage, err := NewCloudstorage(ctx)
	//	exist, err := storage.IsBucketExist(ctx, bucketName)
	//
	IsBucketExist(ctx context.Context, bucketName string) (bool, error)

	// WriteText file to bucket
	//
	//	ctx := context.Background()
	//	storage, err := NewCloudstorage(ctx)
	//	err = storage.AddBucket(ctx, bucketName, "US")
	//	So(err, ShouldBeNil)
	//
	WriteText(ctx context.Context, bucketName, path, txt string) error

	// ReadText file from bucket
	//
	//	ctx := context.Background()
	//	storage, err := NewCloudstorage(ctx)
	//	txt, err := storage.ReadText(ctx, bucketName, path)
	//	So(err, ShouldBeNil)
	//	So(txt, ShouldEqual, "hi")
	//
	ReadText(ctx context.Context, bucketName, path string) (string, error)

	// Delete file from bucket
	//
	//	ctx := context.Background()
	//	storage, err := NewCloudstorage(ctx)
	//	err = storage.Delete(ctx, bucketName, path)
	//
	Delete(ctx context.Context, bucketName, path string) error

	//setCors
}

// CloudstorageImpl is cloudflare implementation
//
type CloudstorageImpl struct {
	Cloudstorage
	client    *storage.Client
	projectID string
}

// NewCloudstorage create Cloudstorage
//
//	ctx := context.Background()
//	storage, err := NewCloudstorage(ctx)
//
func NewCloudstorage(ctx context.Context) (Cloudstorage, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	cred, err := gcp.GlobalCredential(ctx)
	if err != nil {
		return nil, err
	}

	client, err := storage.NewClient(ctx, option.WithCredentials(cred))
	if err != nil {
		return nil, err
	}

	cloudstorage := &CloudstorageImpl{
		client:    client,
		projectID: cred.ProjectID,
	}
	return cloudstorage, nil
}

// AddBucket add cloud storage bucket
//
//	ctx := context.Background()
//	storage, err := NewCloudstorage(ctx)
//	err = storage.AddBucket(ctx, "mock-libsrv.piyuo.com", "US")
//
func (impl *CloudstorageImpl) AddBucket(ctx context.Context, bucketName, location string) error {

	exist, err := impl.IsBucketExist(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exist {
		bucket := impl.client.Bucket(bucketName)
		ctx, cancel := context.WithTimeout(ctx, time.Second*12)
		defer cancel()
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
//	storage, err := NewCloudstorage(ctx)
//	err = storage.RemoveBucket(ctx, "mock-libsrv.piyuo.com")
//
func (impl *CloudstorageImpl) RemoveBucket(ctx context.Context, bucketName string) error {

	exist, err := impl.IsBucketExist(ctx, bucketName)
	if err != nil {
		return err
	}

	if exist {
		bucket := impl.client.Bucket(bucketName)
		ctx, cancel := context.WithTimeout(ctx, time.Second*12)
		defer cancel()
		if err := bucket.Delete(ctx); err != nil {
			return errors.Wrap(err, "failed to remove bucket:"+bucketName)
		}

		log.Info(ctx, here, bucketName+" Bucket deleted")
	}
	return nil
}

// IsBucketExist return true if bucket exist
//
//	bucketName := "mock-libsrv.piyuo.com"
//	ctx := context.Background()
//	storage, err := NewCloudstorage(ctx)
//	exist, err := storage.IsBucketExist(ctx, bucketName)
//
func (impl *CloudstorageImpl) IsBucketExist(ctx context.Context, bucketName string) (bool, error) {

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

// WriteText file to bucket
//
//	ctx := context.Background()
//	storage, err := NewCloudstorage(ctx)
//	err = storage.AddBucket(ctx, bucketName, "US")
//	So(err, ShouldBeNil)
//
func (impl *CloudstorageImpl) WriteText(ctx context.Context, bucketName, path, txt string) error {
	bucket := impl.client.Bucket(bucketName)
	ctx, cancel := context.WithTimeout(ctx, time.Second*12)
	defer cancel()

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
//	storage, err := NewCloudstorage(ctx)
//	txt, err := storage.ReadText(ctx, bucketName, path)
//	So(err, ShouldBeNil)
//	So(txt, ShouldEqual, "hi")
//
func (impl *CloudstorageImpl) ReadText(ctx context.Context, bucketName, path string) (string, error) {
	bucket := impl.client.Bucket(bucketName)
	ctx, cancel := context.WithTimeout(ctx, time.Second*12)
	defer cancel()

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

// Delete file from bucket
//
//	ctx := context.Background()
//	storage, err := NewCloudstorage(ctx)
//	err = storage.Delete(ctx, bucketName, path)
//
func (impl *CloudstorageImpl) Delete(ctx context.Context, bucketName, path string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*12)
	defer cancel()

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
//	storage, err := NewCloudstorage(ctx)
//	err = storage.Delete(ctx, bucketName, path)
//
func (impl *CloudstorageImpl) CleanBucket(ctx context.Context, bucketName string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
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
//
//
func (impl *CloudstorageImpl) RemoveObjects(ctx context.Context, bucket *storage.BucketHandle) (bool, error) {

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
