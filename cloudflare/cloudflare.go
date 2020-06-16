package data

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	file "github.com/piyuo/libsrv/file"
	"github.com/pkg/errors"
)

// Cloudflare is dns toolkit for implement cloudflare api
//
type Cloudflare interface {

	// AddSubDomain add sub domain cname record
	//
	AddSubDomain(ctx context.Context, domainName string) error

	// RemoveSubDomain remove sub domain cname record
	//
	RemoveSubDomain(ctx context.Context, domainName string) error

	// IsSubDomainExist return true if sub domain exist
	//
	IsSubDomainExist(ctx context.Context, domainName string) (bool, error)
}

// CloudflareImpl is cloudflare implementation
//
type CloudflareImpl struct {
	Cloudflare
	zone    string
	account string
	token   string
}

// NewCloudflare create Cloudflare
//
func NewCloudflare(ctx context.Context) (Cloudflare, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	keyPath := "keys/cloudflare.key"
	currentDir, _ := os.Getwd()
	keyDir := path.Join(currentDir, "../../"+keyPath)
	keyFile, err := file.Open(keyDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get "+keyPath)
	}
	defer keyFile.Close()

	json, err := keyFile.JSON()
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert "+keyPath+" to JSON")
	}

	cloudflare := &CloudflareImpl{
		zone:    json["zone"].(string),
		account: json["account"].(string),
		token:   json["token"].(string),
	}
	return cloudflare, nil
}

// AddSubDomain add sub domain cname record
//
//	fmt.Println(f.JSON()["users"])
//
func (impl *CloudflareImpl) AddSubDomain(ctx context.Context, domainName string) error {

	url := "https://api.cloudflare.com/client/v4/zones/" + impl.zone + "/dns_records"
	var jsonStr = []byte(`{"type":"CNAME","name":"` + domainName + `","content":"ghs.googlehosted.com","ttl":1,"priority":10,"proxied":false}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer "+impl.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to add sub domain: "+domainName)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return nil
}
