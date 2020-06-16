package data

import (
	"context"
	"time"

	"cloud.google.com/go/storage"
	"github.com/piyuo/libsrv/log"
	"github.com/piyuo/libsrv/secure/gcp"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

const here = "cloudstorage"

// Cloudstorage is google cloud storage toolkit
//
type Cloudstorage interface {

	// AddBucket add cloud storage bucket
	//
	AddBucket(ctx context.Context, bucketName string) error

	// RemoveBucket remove cloud storage bucket
	//
	RemoveBucket(ctx context.Context, bucketName string) error

	// IsBucketExist return true if bucket exist
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

	return &CloudstorageImpl{
		client:    client,
		projectID: cred.ProjectID,
	}, nil
}

// AddBucket add cloud storage bucket
//
//	fmt.Println(f.JSON()["users"])
//
func (impl *CloudstorageImpl) AddBucket(ctx context.Context, bucketName string) error {

	// Creates a Bucket instance.
	bucket := impl.client.Bucket(bucketName)

	// Creates the new bucket.
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := bucket.Create(ctx, impl.projectID, &storage.BucketAttrs{
		StorageClass: "COLDLINE",
		Location:     "asia",
	}); err != nil {
		return errors.Wrap(err, "failed to create bucket "+bucketName)
	}

	log.Info(ctx, here, bucketName+" Bucket created")
	return nil
}
