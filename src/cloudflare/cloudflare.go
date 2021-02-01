package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	key "github.com/piyuo/libsrv/src/key"
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

	// AddTxtRecord add TXT record to dns
	//
	//	ctx := context.Background()
	//	cflare, err := NewCloudflare(ctx)
	//	exist, err := cflare.IsDomainExist(ctx, domainName)
	//
	AddTxtRecord(ctx context.Context, domainName, txt string) error

	// RemoveTxtRecord removeTXT record from dns
	//
	//	ctx := context.Background()
	//	cflare, err := NewCloudflare(ctx)
	//	err = cflare.RemoveTxtRecord(ctx, domainName, txt)
	//
	RemoveTxtRecord(ctx context.Context, domainName string) error

	// IsTxtRecordExist return true if txt record exist
	//
	//	ctx := context.Background()
	//	cflare, err := NewCloudflare(ctx)
	//	exist, err := cflare.IsTxtRecordExist(ctx, domainName, txt)
	//
	IsTxtRecordExist(ctx context.Context, domainName string) (bool, error)
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

	json, err := key.JSON("cloudflare.json")
	if err != nil {
		return nil, err
	}
	cloudflare := &CloudflareImpl{
		zone:    json["zone"].(string),
		account: json["account"].(string),
		token:   json["token"].(string),
	}
	return cloudflare, nil
}

// getDomainID return domain id if domain exist otherwise return empty string
//
//	id, err := impl.getDomainID(ctx, domainName)
//
func (impl *CloudflareImpl) getRecord(ctx context.Context, domainName, recType, content string) (string, error) {
	url := impl.apiURL() + "?type=" + recType + "&name=" + url.QueryEscape(domainName)
	if content != "" {
		url = url + "&content=" + content
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to get record:"+domainName)
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

	body, _ := ioutil.ReadAll(resp.Body)
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode json")
	}

	success := response["success"].(bool)
	if !success {
		message := ""
		dnsErrors := response["errors"].([]interface{})
		for _, dnsErrorRecord := range dnsErrors {
			dnsError := dnsErrorRecord.(map[string]interface{})
			message += dnsError["message"].(string)
		}
		return nil, errors.New(message)
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

	var requestJSON = []byte(`{"type":"CNAME","name":"` + domainName + `","content":"ghs.googlehosted.com","ttl":1,"priority":10,"proxied":` + proxy + `}`)
	req, err := http.NewRequest("POST", impl.apiURL(), bytes.NewBuffer(requestJSON))
	if err != nil {
		return errors.Wrap(err, "failed to new add domain:"+domainName)
	}

	_, err = impl.dnsRequest(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return err
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
	id, err := impl.getRecord(ctx, domainName, "CNAME", "ghs.googlehosted.com")
	if err != nil {
		return err
	}
	if id != "" {
		url := impl.apiURL() + "/" + id
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return errors.Wrap(err, "failed to remove domain:"+domainName)
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

	id, err := impl.getRecord(ctx, domainName, "CNAME", "ghs.googlehosted.com")
	if err != nil {
		return false, err
	}
	return id != "", nil
}

// AddTxtRecord add TXT record to dns
//
//	ctx := context.Background()
//	cflare, err := NewCloudflare(ctx)
//	err = cflare.AddTxtRecord(ctx, domainName, txt)
//
func (impl *CloudflareImpl) AddTxtRecord(ctx context.Context, domainName, txt string) error {
	var requestJSON = []byte(`{"type":"TXT","name":"` + domainName + `","content":"` + txt + `"}`)
	req, err := http.NewRequest("POST", impl.apiURL(), bytes.NewBuffer(requestJSON))
	if err != nil {
		return errors.Wrap(err, "failed to add txt record:"+domainName+" context:"+txt)
	}
	_, err = impl.dnsRequest(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return err
	}
	return nil
}

// RemoveTxtRecord remove txt record from dns
//
//	ctx := context.Background()
//	cflare, err := NewCloudflare(ctx)
//	err = cflare.RemoveTxtRecord(ctx, domainName, txt)
//
func (impl *CloudflareImpl) RemoveTxtRecord(ctx context.Context, domainName string) error {
	id, err := impl.getRecord(ctx, domainName, "TXT", "")
	if err != nil {
		return err
	}
	if id != "" {
		url := impl.apiURL() + "/" + id
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return errors.Wrap(err, "failed to remove txt record:"+domainName)
		}
		_, err = impl.dnsRequest(ctx, req)
		if err != nil {
			return err
		}
	}
	return nil
}

// IsTxtRecordExist return true if txt record exist
//
//	ctx := context.Background()
//	cflare, err := NewCloudflare(ctx)
//	exist, err = cflare.IsTxtRecordExist(ctx, domainName, txt)
//
func (impl *CloudflareImpl) IsTxtRecordExist(ctx context.Context, domainName string) (bool, error) {
	id, err := impl.getRecord(ctx, domainName, "TXT", "")
	if err != nil {
		return false, err
	}
	return id != "", nil
}
