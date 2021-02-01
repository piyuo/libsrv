package cloudflare

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCloudflareCredential(t *testing.T) {
	assert := assert.New(t)
	zone, token, err := credential()
	assert.Nil(err)
	assert.NotEmpty(zone)
	assert.NotEmpty(token)
}

func TestCloudflareSendDNSRequest(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	resp, err := sendDNSRequest(ctx, "GET", "?type=CNAME", nil)
	assert.Nil(err)
	assert.NotEmpty(resp["result"])
}

func TestCloudflareGetDNSRecordID(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	id, err := getDNSRecordID(ctx, "piyuo.com", "CNAME", "")
	assert.Nil(err)
	assert.NotEmpty(id)
}

func TestDomain(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	subDomain := "mock-libsrv"
	domainName := subDomain + ".piyuo.com"

	//remove sample domain
	RemoveDomain(ctx, domainName)

	exist, err := IsDomainExist(ctx, domainName)
	assert.Nil(err)
	assert.False(exist)

	err = AddDomain(ctx, domainName, false)
	assert.Nil(err)

	exist, err = IsDomainExist(ctx, domainName)
	assert.Nil(err)
	assert.True(exist)

	// add domain that already exist should not error
	err = AddDomain(ctx, domainName, false)
	assert.Nil(err)

	err = RemoveDomain(ctx, domainName)
	assert.Nil(err)

	exist, err = IsDomainExist(ctx, domainName)
	assert.Nil(err)
	assert.False(exist)

	// remove domain second time should not error
	err = RemoveDomain(ctx, domainName)
	assert.Nil(err)

	TestMode = true
	err = AddDomain(ctx, domainName, false)
	assert.Nil(err)

	exist, err = IsDomainExist(ctx, domainName)
	assert.Nil(err)
	assert.True(exist)

	err = RemoveDomain(ctx, domainName)
	assert.Nil(err)
	TestMode = false

}

func TestTxtRecord(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	subDomain := "mock-libsrv"
	domainName := subDomain + ".piyuo.com"
	txt := "hi"
	//remove sample record
	RemoveTxtRecord(ctx, domainName)

	exist, err := IsTxtRecordExist(ctx, domainName)
	assert.Nil(err)
	assert.False(exist)

	err = AddTxtRecord(ctx, domainName, txt)
	assert.Nil(err)

	exist, err = IsTxtRecordExist(ctx, domainName)
	assert.Nil(err)
	assert.True(exist)

	// add txt record that already exist should not error
	err = AddTxtRecord(ctx, domainName, txt)
	assert.Nil(err)

	err = RemoveTxtRecord(ctx, domainName)
	assert.Nil(err)

	exist, err = IsTxtRecordExist(ctx, domainName)
	assert.Nil(err)
	assert.False(exist)

	// remove txt record second time should not error
	err = RemoveTxtRecord(ctx, domainName)
	assert.Nil(err)

	TestMode = true
	err = AddTxtRecord(ctx, domainName, txt)
	assert.Nil(err)
	err = RemoveTxtRecord(ctx, domainName)
	assert.Nil(err)
	exist, err = IsTxtRecordExist(ctx, domainName)
	assert.Nil(err)
	assert.True(exist)
	TestMode = false

}
