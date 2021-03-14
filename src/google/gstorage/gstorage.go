package gstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/storage"
	"github.com/piyuo/libsrv/src/log"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
)

// Gstorage is google cloud storage toolkit
//
//	before use this toolkit you need verify domain owner ship, please add service account in gclould.key to  https://www.google.com/webmasters/verification
//
type Gstorage interface {

	// CreateBucket add cloud storage bucket
	//
	//	storage, err := New(ctx)
	//	err = storage.CreateBucket(ctx, "my-bucket","us-central1","region")
	//
	CreateBucket(ctx context.Context, bucketName, location, locationType string) error

	// DeleteBucket remove cloud storage bucket
	//
	//	storage, err := New(ctx)
	//	err = storage.DeleteBucket(ctx, "my-bucket")
	//
	DeleteBucket(ctx context.Context, bucketName string) error

	// PublicBucket make bucket public
	//
	//	storage, err := New(ctx)
	//	err = storage.PublicBucket(ctx, "my-bucket")
	//
	PublicBucket(ctx context.Context, bucketName string) error

	// SetPageAndCORS set bucket CORS configuration
	//
	//	storage, err := New(ctx)
	//	err = storage.SetPageAndCORS(ctx, "my-bucket","some-origin.com")
	//
	SetPageAndCORS(ctx context.Context, bucketName, originDomain string) error

	// CleanBucket remove all file in bucket
	//
	//	storage, err := New(ctx)
	//	err = storage.RemoveBucket(ctx, "my-bucket")
	//
	CleanBucket(ctx context.Context, bucketName string) error

	// IsBucketExists return true if bucket exist
	//
	//	bucketName := "my-bucket"
	//	storage, err := New(ctx)
	//	exist, err := storage.IsBucketExists(ctx, bucketName)
	//
	IsBucketExists(ctx context.Context, bucketName string) (bool, error)

	// IsFileExists return true if file exist
	//
	//	storage, err := New(ctx)
	//	found,err = storage.IsFileExists(ctx, bucketName,"fileName")
	//
	IsFileExists(ctx context.Context, bucketName, fileName string) (bool, error)

	// ListFiles list bucket files base on prefix and delimiters
	//
	ListFiles(ctx context.Context, bucketName, prefix, delim string) ([]string, error)

	// DeleteFiles delete files in prefix
	//
	//	storage, err := New(ctx)
	//	err = storage.DeleteFiles(ctx, bucketName, "assets")
	//
	DeleteFiles(ctx context.Context, bucketName, prefix string) error

	// WriteText write text file to bucket
	//
	//	storage, err := New(ctx)
	//	err = storage.WriteText(ctx, bucketName, "a/b.txt")
	//
	WriteText(ctx context.Context, bucketName, filename, txt string) error

	// WriteJSON write json file to bucket
	//
	//	storage, err := New(ctx)
	//	err = storage.WriteJSON(ctx, bucketName, "a/b.json")
	//
	WriteJSON(ctx context.Context, bucketName, filename string, context map[string]interface{}) error

	// ReadText file from bucket
	//
	//	storage, err := New(ctx)
	//	txt, err := storage.ReadText(ctx, bucketName, "a/b.txt")
	//
	ReadText(ctx context.Context, bucketName, path string) (string, error)

	// ReadJSON read json file from bucket
	//
	//	storage, err := New(ctx)
	//	txt, err := storage.ReadJSON(ctx, bucketName, "a/b.json")
	//
	ReadJSON(ctx context.Context, bucketName, path string) (map[string]interface{}, error)

	// DeleteFile file from bucket
	//
	//	storage, err := New(ctx)
	//	err = storage.DeleteFile(ctx, bucketName, path)
	//
	DeleteFile(ctx context.Context, bucketName, path string) error

	// Sync bucket to dir, remove bucket files not exist in dir
	//
	//	shouldDeleteFile := func(filename string) (bool, error) {
	//		return true, nil
	//	}
	//	storage, err := New(ctx)
	//	err = storage.Sync(ctx, bucketName, dir)
	//
	Sync(ctx context.Context, bucketName string, shouldDeleteFile ShouldDeleteFile) error
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

// CreateBucket add cloud storage bucket
//
//	storage, err := New(ctx)
//	err = storage.CreateBucket(ctx, "my-bucket","us-central1","region")
//
func (impl *Implementation) CreateBucket(ctx context.Context, bucketName, location, locationType string) error {

	bucket := impl.client.Bucket(bucketName)
	if err := bucket.Create(ctx, impl.projectID, &storage.BucketAttrs{
		Location:     "us-central1",
		LocationType: "region",
	}); err != nil {
		return errors.Wrap(err, "failed to add bucket:"+bucketName)
	}
	log.Info(ctx, bucketName+" Bucket created")
	return nil
}

// DeleteBucket remove cloud storage bucket
//
//	storage, err := New(ctx)
//	err = storage.DeleteBucket(ctx, "my-bucket")
//
func (impl *Implementation) DeleteBucket(ctx context.Context, bucketName string) error {

	bucket := impl.client.Bucket(bucketName)

	if err := impl.CleanBucket(ctx, bucketName); err != nil {
		return errors.Wrap(err, "failed to clean bucket before remove:"+bucketName)
	}

	if err := bucket.Delete(ctx); err != nil {
		return errors.Wrap(err, "failed to remove bucket:"+bucketName)
	}
	log.Info(ctx, bucketName+" Bucket deleted")
	return nil
}

// PublicBucket make bucket public
//
//	storage, err := New(ctx)
//	err = storage.PublicBucket(ctx, "my-bucket")
//
func (impl *Implementation) PublicBucket(ctx context.Context, bucketName string) error {
	bucket := impl.client.Bucket(bucketName)
	policy, err := bucket.IAM().V3().Policy(ctx)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to get Bucket(%q).IAM().V3().Policy", bucketName))
	}
	role := "roles/storage.objectViewer"
	policy.Bindings = append(policy.Bindings, &iampb.Binding{
		Role:    role,
		Members: []string{iam.AllUsers},
	})
	if err := bucket.IAM().V3().SetPolicy(ctx, policy); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to set Bucket(%q).IAM().SetPolicy", bucketName))
	}
	return nil
}

// SetPageAndCORS set bucket WEB & CORS configuration, let it be static website
//
//	storage, err := New(ctx)
//	err = storage.SetPageAndCORS(ctx, "my-bucket","some-origin.com")
//
func (impl *Implementation) SetPageAndCORS(ctx context.Context, bucketName, originDomain string) error {
	bucket := impl.client.Bucket(bucketName)
	bucketAttrsToUpdate := storage.BucketAttrsToUpdate{
		Website: &storage.BucketWebsite{
			MainPageSuffix: "index.html",
			NotFoundPage:   "404.html",
		},
		CORS: []storage.CORS{
			{
				MaxAge:  time.Second * 86400, //Maximum number of seconds the results can be cached, we use 24 hours
				Methods: []string{"GET", "POST"},
				Origins: []string{originDomain}, //[]string{"*"},
			}},
	}
	if _, err := bucket.Update(ctx, bucketAttrsToUpdate); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to exec Bucket(%q).Update", bucketName))
	}
	return nil
}

// IsBucketExists return true if bucket exist
//
//	bucketName := "my-bucket
//	storage, err := New(ctx)
//	exist, err := storage.IsBucketExists(ctx, bucketName)
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
//	storage, err := New(ctx)
//	found,err = storage.IsFileExists(ctx, bucketName, "fileName")
//
func (impl *Implementation) IsFileExists(ctx context.Context, bucketName, path string) (bool, error) {
	bucket := impl.client.Bucket(bucketName)
	query := &storage.Query{}
	it := bucket.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, err
		}
		if attrs.Name == path {
			return true, nil
		}
	}
	return false, nil
}

// WriteText file to bucket
//
//	err = storage.WriteText(ctx, bucketName, "a/b.txt")
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
//	txt, err := storage.ReadText(ctx, bucketName, "a/b.txt")
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

// WriteJSON write json file to bucket
//
//	err = storage.WriteJSON(ctx, bucketName, "a/b.json")
//
func (impl *Implementation) WriteJSON(ctx context.Context, bucketName, path string, content map[string]interface{}) error {
	bytes, err := json.Marshal(content)
	if err != nil {
		return errors.Wrap(err, "failed to marshal version control to json")
	}
	if err = impl.WriteText(ctx, bucketName, path, string(bytes)); err != nil {
		return errors.Wrap(err, "failed to write verstion control to bucket")
	}
	return nil
}

// ReadJSON read json file from bucket. it will return empty json if file not exist
//
//	txt, err := storage.ReadJSON(ctx, bucketName, "a/b.json")
//
func (impl *Implementation) ReadJSON(ctx context.Context, bucketName, path string) (map[string]interface{}, error) {
	var control map[string]interface{}
	text, err := impl.ReadText(ctx, bucketName, path)
	if err == nil {
		json.Unmarshal([]byte(text), &control)
	}
	if control == nil {
		control = map[string]interface{}{}
	}
	return control, nil
}

// ListFiles list all files
//
//	files, err := storage.ListFiles(ctx, bucketName, "", "")
//
func (impl *Implementation) ListFiles(ctx context.Context, bucketName, prefix, delim string) ([]string, error) {
	files := []string{}
	bucket := impl.client.Bucket(bucketName)
	query := &storage.Query{Prefix: prefix, Delimiter: delim}
	query.SetAttrSelection([]string{"Name"})

	it := bucket.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		files = append(files, attrs.Name)
	}
	return files, nil
}

// DeleteFiles delete files in prefix
//
//	storage, err := New(ctx)
//	err = storage.DeleteFiles(ctx, bucketName, "assets")
//
func (impl *Implementation) DeleteFiles(ctx context.Context, bucketName, prefix string) error {
	query := &storage.Query{Prefix: prefix}
	query.SetAttrSelection([]string{"Name"})
	it := impl.client.Bucket(bucketName).Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		obj := impl.client.Bucket(bucketName).Object(attrs.Name)
		if err := obj.Delete(ctx); err != nil {
			return err
		}
		log.Debug(ctx, fmt.Sprintf("delete object:%v", attrs.Name))
	}
	return nil
}

// DeleteFile file from bucket
//
//	storage, err := New(ctx)
//	err = storage.DeleteFile(ctx, bucketName, path)
//
func (impl *Implementation) DeleteFile(ctx context.Context, bucketName, path string) error {

	o := impl.client.Bucket(bucketName).Object(path)
	if err := o.Delete(ctx); err != nil {
		return err
	}
	log.Debug(ctx, fmt.Sprintf("delete object:%v", path))
	return nil
}

// CleanBucket remove all files in bucket, return true if still have file in bucket
// timeout in ms
//
//	storage, err := New(ctx)
//	err = storage.Delete(ctx, bucketName, path)
//
func (impl *Implementation) CleanBucket(ctx context.Context, bucketName string) error {
	bucket := impl.client.Bucket(bucketName)
	for {
		result, err := impl.removeObjects(ctx, bucket)
		if err != nil {
			return err
		}
		if result == true {
			return nil
		}
	}
}

// removeObjects remove objects max 1000, return true if object all deleted
//
func (impl *Implementation) removeObjects(ctx context.Context, bucket *storage.BucketHandle) (bool, error) {
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
		if err := bucket.Object(attrs.Name).Delete(ctx); err != nil {
			return false, err
		}
		log.Debug(ctx, fmt.Sprintf("clean object:%v, i=%v", attrs.Name, i))
		i++
		if i >= 1000 {
			return false, nil
		}
	}
	return true, nil
}

// ShouldDeleteFile use in SyncDir to determine a file should be delete
//
type ShouldDeleteFile func(filename string) (bool, error)

// Sync bucket to dir, remove bucket files not exist in dir
//
//	shouldDeleteFile := func(filename string) (bool, error) {
//		return true, nil
//	}
//	storage, err := New(ctx)
//	err = storage.Sync(ctx, bucketName, dir)
//
func (impl *Implementation) Sync(ctx context.Context, bucketName string, shouldDeleteFile ShouldDeleteFile) error {

	bucket := impl.client.Bucket(bucketName)
	query := &storage.Query{}
	query.SetAttrSelection([]string{"Name"})

	it := bucket.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		delete, err := shouldDeleteFile(attrs.Name)
		if err != nil {
			return err
		}
		if delete {
			if err := bucket.Object(attrs.Name).Delete(ctx); err != nil {
				return err
			}
			log.Debug(ctx, fmt.Sprintf("sync delete object:%v", attrs.Name))
		}
	}
	return nil
}
