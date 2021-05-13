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

	key "github.com/piyuo/libsrv/key"
	"github.com/pkg/errors"
)

// Mock define key test flag
//
type Mock int8

const (
	// MockNoError let function return nil
	//
	MockNoError Mock = iota

	// MockError let function error
	//
	MockError

	// CnameExists let IsCNAMEExists return exists
	//
	MockCnameNotExists
)

//	credential return cloudflare credential zon and token
//
//	zone,token, err = impl.credential()
//
func credential() (string, string, error) {
	json, err := key.JSON("cloudflare.json")
	if err != nil {
		return "", "", errors.Wrap(err, "get key from cloudflare.json")
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
		return nil, errors.Wrap(err, "get credential")
	}

	url := "https://api.cloudflare.com/client/v4/zones/" + zone + "/dns_records" + query
	req, err := http.NewRequest(method, url, reqestBody)
	if err != nil {
		return nil, errors.Wrap(err, "http request")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	//	ctx, cancel := context.WithTimeout(ctx, time.Second*15) // cloud flare dns call must completed in 15 seconds
	//	defer cancel()

	client := &http.Client{
		Timeout: time.Duration(time.Second * 15),
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "dns request")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.Wrap(err, "decode json")
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

// CreateCloudRunCNAME add domain cname record point to google cloud run backend, it will create not proxied CNAME and after manually add cloud run mapping you should manually set proxied=true on cloudflare
//
//	err = CreateCloudRunCNAME(ctx, "my.piyuo.com")
//
func CreateCloudRunCNAME(ctx context.Context, domainName string) error {
	return CreateCNAME(ctx, domainName, "ghs.googlehosted.com", false)
}

// CreateStorageCNAME add domain cname record point to google storage backend
//
//	err = CreateStorageCNAME(ctx, "my.piyuo.com")
//
func CreateStorageCNAME(ctx context.Context, domainName string) error {
	return CreateCNAME(ctx, domainName, "c.storage.googleapis.com", true)
}

// CreateCNAME create domain CNAME record
//
//	err = AddCNAME(ctx, "my.piyuo.com", false)
//
func CreateCNAME(ctx context.Context, domainName, target string, proxied bool) error {
	if ctx.Value(MockNoError) != nil {
		return nil
	}
	if ctx.Value(MockError) != nil {
		return errors.New("")
	}

	proxy := "false"
	if proxied {
		proxy = "true"
	}

	var requestJSON = []byte(`{"type":"CNAME","name":"` + domainName + `","content":"` + target + `","ttl":1,"priority":10,"proxied":` + proxy + `}`)
	_, err := sendDNSRequest(ctx, "POST", "", bytes.NewBuffer(requestJSON))
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return err
	}
	return nil
}

// DeleteCNAME delete cname record, return no error if domain name not exists
//
//	err = DeleteCNAME(ctx, "my.piyuo.com")
//
func DeleteCNAME(ctx context.Context, domainName string) error {
	if ctx.Value(MockNoError) != nil {
		return nil
	}
	if ctx.Value(MockError) != nil {
		return errors.New("")
	}

	id, err := getDNSRecordID(ctx, domainName, "CNAME", "")
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

// IsCNAMEExists return true if CNAME exist
//
//	exist, err := IsCNAMEExists(ctx, "my.piyuo.com")
//
func IsCNAMEExists(ctx context.Context, domainName string) (bool, error) {
	if ctx.Value(MockCnameNotExists) != nil {
		return false, nil
	}
	if ctx.Value(MockNoError) != nil {
		return true, nil
	}
	if ctx.Value(MockError) != nil {
		return false, errors.New("")
	}

	id, err := getDNSRecordID(ctx, domainName, "CNAME", "")
	if err != nil {
		return false, err
	}
	return id != "", nil
}

// CreateTXT add TXT record to dns
//
//	err = cflare.CreateTXT(ctx, "my.piyuo.com", txt)
//
func CreateTXT(ctx context.Context, domainName, txt string) error {
	if ctx.Value(MockNoError) != nil {
		return nil
	}
	if ctx.Value(MockError) != nil {
		return errors.New("")
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

// RemoveTXT remove txt record from dns
//
//	err = RemoveTXT(ctx, "my.piyuo.com", txt)
//
func RemoveTXT(ctx context.Context, domainName string) error {
	if ctx.Value(MockNoError) != nil {
		return nil
	}
	if ctx.Value(MockError) != nil {
		return errors.New("")
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

// IsTXTExists return true if txt record exist
//
//	exist, err = IsTXTExists(ctx, "my.piyuo.com", txt)
//
func IsTXTExists(ctx context.Context, domainName string) (bool, error) {
	if ctx.Value(MockNoError) != nil {
		return true, nil
	}
	if ctx.Value(MockError) != nil {
		return false, errors.New("")
	}

	id, err := getDNSRecordID(ctx, domainName, "TXT", "")
	if err != nil {
		return false, err
	}
	return id != "", nil
}
