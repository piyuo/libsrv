package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	key "github.com/piyuo/libsrv/src/key"
	"github.com/pkg/errors"
)

// testMode set to true will let every function run success
//
var testMode = false

// EnableTestMode set to true will let every function run success
//
func EnableTestMode(enabled bool) {
	testMode = enabled
}

//	credential return cloudflare credential zon and token
//
//	zone,token, err = impl.credential()
//
func credential() (string, string, error) {
	json, err := key.JSON("cloudflare.json")
	if err != nil {
		return "", "", errors.Wrap(err, "failed to get key from JSON:cloudflare.json")
	}
	return json["zone"].(string), json["token"].(string), nil
}

//	sendDNSRequest add authorization to request and check response is success
//
//	response, err = sendDNSRequest(ctx, req)
//
func sendDNSRequest(ctx context.Context, method, query string, reqestBody io.Reader) (map[string]interface{}, error) {
	zone, token, err := credential()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get credential")
	}

	url := "https://api.cloudflare.com/client/v4/zones/" + zone + "/dns_records" + query
	req, err := http.NewRequest(method, url, reqestBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to new http request")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(ctx, time.Second*15) // cloud flare dns call must completed in 15 seconds
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

// getDNSRecordID return dns record id if domain exist otherwise return empty string
//
//	id, err := getDNSRecordID(ctx, "piyuo.com", "CNAME", "")
//
func getDNSRecordID(ctx context.Context, domainName, recType, content string) (string, error) {
	query := "?type=" + recType + "&name=" + url.QueryEscape(domainName)
	if content != "" {
		query = query + "&content=" + content
	}

	response, err := sendDNSRequest(ctx, "GET", query, nil)
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

// AddDomain add domain cname record
//
//	err = AddDomain(ctx, domainName, false)
//
func AddDomain(ctx context.Context, domainName string, proxied bool) error {
	if testMode {
		return nil
	}

	proxy := "false"
	if proxied {
		proxy = "true"
	}

	var requestJSON = []byte(`{"type":"CNAME","name":"` + domainName + `","content":"ghs.googlehosted.com","ttl":1,"priority":10,"proxied":` + proxy + `}`)
	_, err := sendDNSRequest(ctx, "POST", "", bytes.NewBuffer(requestJSON))
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
//	err = RemoveDomain(ctx, domainName)
//
func RemoveDomain(ctx context.Context, domainName string) error {
	if testMode {
		return nil
	}

	id, err := getDNSRecordID(ctx, domainName, "CNAME", "ghs.googlehosted.com")
	if err != nil {
		return err
	}
	if id != "" {
		query := "/" + id
		_, err = sendDNSRequest(ctx, "DELETE", query, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// IsDomainExist return true if domain exist
//
//	exist, err := IsDomainExist(ctx, domainName)
//
func IsDomainExist(ctx context.Context, domainName string) (bool, error) {
	if testMode {
		return true, nil
	}

	id, err := getDNSRecordID(ctx, domainName, "CNAME", "ghs.googlehosted.com")
	if err != nil {
		return false, err
	}
	return id != "", nil
}

// AddTxtRecord add TXT record to dns
//
//	err = cflare.AddTxtRecord(ctx, domainName, txt)
//
func AddTxtRecord(ctx context.Context, domainName, txt string) error {
	if testMode {
		return nil
	}

	var requestJSON = []byte(`{"type":"TXT","name":"` + domainName + `","content":"` + txt + `"}`)
	_, err := sendDNSRequest(ctx, "POST", "", bytes.NewBuffer(requestJSON))
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
//	err = RemoveTxtRecord(ctx, domainName, txt)
//
func RemoveTxtRecord(ctx context.Context, domainName string) error {
	if testMode {
		return nil
	}

	id, err := getDNSRecordID(ctx, domainName, "TXT", "")
	if err != nil {
		return err
	}
	if id != "" {
		query := "/" + id
		_, err = sendDNSRequest(ctx, "DELETE", query, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// IsTxtRecordExist return true if txt record exist
//
//	exist, err = IsTxtRecordExist(ctx, domainName, txt)
//
func IsTxtRecordExist(ctx context.Context, domainName string) (bool, error) {
	if testMode {
		return true, nil
	}

	id, err := getDNSRecordID(ctx, domainName, "TXT", "")
	if err != nil {
		return false, err
	}
	return id != "", nil
}
