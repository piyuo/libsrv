package data

import (
	"context"
	"time"

	"cloud.google.com/go/storage"
	"github.com/piyuo/libsrv/log"
	"github.com/piyuo/libsrv/secure/gcp"
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

	// IsBucketExist return true if bucket exist
	//
	//	bucketName := "mock-libsrv.piyuo.com"
	//	ctx := context.Background()
	//	storage, err := NewCloudstorage(ctx)
	//	exist, err := storage.IsBucketExist(ctx, bucketName)
	//
	IsBucketExist(ctx context.Context, bucketName string) (bool, error)
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
