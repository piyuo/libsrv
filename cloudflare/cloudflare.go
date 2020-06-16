package data

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	file "github.com/piyuo/libsrv/file"
	"github.com/pkg/errors"
)

// Cloudflare is dns toolkit for implement cloudflare api
//
type Cloudflare interface {

	// AddDomain add domain cname record
	//
	//	ctx := context.Background()
	//	cflare, err := NewCloudflare(ctx)
	//	err = cflare.AddDomain(ctx, domainName, false)
	//
	AddDomain(ctx context.Context, domainName string, proxied bool) error

	// RemoveDomain remove sub domain cname record
	//
	//	ctx := context.Background()
	//	cflare, err := NewCloudflare(ctx)
	//	err = cflare.RemoveDomain(ctx, domainName)
	//
	RemoveDomain(ctx context.Context, domainName string) error

	// IsDomainExist return true if domain exist
	//
	//	ctx := context.Background()
	//	cflare, err := NewCloudflare(ctx)
	//	exist, err := cflare.IsDomainExist(ctx, domainName)
	//
	IsDomainExist(ctx context.Context, domainName string) (bool, error)
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
//	cflare, err := NewCloudflare(context.Background())
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

//	dnsRequest add authorization to request and check response is success
//
//	_, err = impl.dnsRequest(ctx, req)
//
func (impl *CloudflareImpl) dnsRequest(ctx context.Context, req *http.Request) (map[string]interface{}, error) {

	req.Header.Set("Authorization", "Bearer "+impl.token)
	req.Header.Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(ctx, time.Second*12)
	defer cancel()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make dns request")
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode json")
	}

	success := response["success"].(bool)
	if !success {
		return nil, errors.New(response["message"].(string))
	}

	return response, nil
}

// apiURL return cloudflare api url
//
//	impl.apiURL()
//
func (impl *CloudflareImpl) apiURL() string {
	return "https://api.cloudflare.com/client/v4/zones/" + impl.zone + "/dns_records"
}

// AddDomain add domain cname record
//
//	ctx := context.Background()
//	cflare, err := NewCloudflare(ctx)
//	err = cflare.AddDomain(ctx, domainName, false)
//
func (impl *CloudflareImpl) AddDomain(ctx context.Context, domainName string, proxied bool) error {

	proxy := "false"
	if proxied {
		proxy = "true"
	}
	exist, err := impl.IsDomainExist(ctx, domainName)
	if err != nil {
		return nil
	}
	if !exist {
		var requestJSON = []byte(`{"type":"CNAME","name":"` + domainName + `","content":"ghs.googlehosted.com","ttl":1,"priority":10,"proxied":` + proxy + `}`)
		req, err := http.NewRequest("POST", impl.apiURL(), bytes.NewBuffer(requestJSON))
		if err != nil {
			return errors.Wrap(err, "failed to new request:"+domainName)
		}

		_, err = impl.dnsRequest(ctx, req)
		if err != nil {
			return err
		}

	}
	return nil
}

// RemoveDomain remove sub domain cname record
//
//	ctx := context.Background()
//	cflare, err := NewCloudflare(ctx)
//	err = cflare.RemoveDomain(ctx, domainName)
//
func (impl *CloudflareImpl) RemoveDomain(ctx context.Context, domainName string) error {
	id, err := impl.getDomainID(ctx, domainName)
	if err != nil {
		return err
	}
	if id != "" {
		url := impl.apiURL() + "/" + id
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return errors.Wrap(err, "failed to new request:"+domainName)
		}
		_, err = impl.dnsRequest(ctx, req)
		if err != nil {
			return err
		}
	}
	return nil
}

// IsDomainExist return true if domain exist
//
//	ctx := context.Background()
//	cflare, err := NewCloudflare(ctx)
//	exist, err := cflare.IsDomainExist(ctx, domainName)
//
func (impl *CloudflareImpl) IsDomainExist(ctx context.Context, domainName string) (bool, error) {

	id, err := impl.getDomainID(ctx, domainName)
	if err != nil {
		return false, err
	}
	return id != "", nil
}

// getDomainID return domain id if domain exist otherwise return empty string
//
//	id, err := impl.getDomainID(ctx, domainName)
//
func (impl *CloudflareImpl) getDomainID(ctx context.Context, domainName string) (string, error) {

	url := impl.apiURL() + "?type=CNAME&content=ghs.googlehosted.com&name=" + url.QueryEscape(domainName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to new request:"+domainName)
	}

	response, err := impl.dnsRequest(ctx, req)
	if err != nil {
		return "", err
	}

	result := response["result"].([]interface{})
	if len(result) == 0 {
		return "", nil
	}
	rec := result[0].(map[string]interface{})
	return rec["id"].(string), nil
}
