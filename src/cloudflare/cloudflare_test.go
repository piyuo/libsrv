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

func TestCloudflareCNAME(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	subDomain := "mock-libsrv"
	domainName := subDomain + ".piyuo.com"

	//remove sample domain
	DeleteCNAME(ctx, domainName)

	exist, err := IsCNAMEExists(ctx, domainName)
	assert.Nil(err)
	assert.False(exist)

	err = CreateCNAME(ctx, domainName, "ghs.googlehosted.com", false)
	assert.Nil(err)

	exist, err = IsCNAMEExists(ctx, domainName)
	assert.Nil(err)
	assert.True(exist)

	// add domain that already exist should not error
	err = CreateCNAME(ctx, domainName, "ghs.googlehosted.com", false)
	assert.Nil(err)

	err = DeleteCNAME(ctx, domainName)
	assert.Nil(err)

	exist, err = IsCNAMEExists(ctx, domainName)
	assert.Nil(err)
	assert.False(exist)

	// remove domain second time should not error
	err = DeleteCNAME(ctx, domainName)
	assert.Nil(err)

	// cloud run
	err = CreateCloudRunCNAME(ctx, domainName)
	assert.Nil(err)
	exist, err = IsCNAMEExists(ctx, domainName)
	assert.Nil(err)
	assert.True(exist)
	err = DeleteCNAME(ctx, domainName)
	assert.Nil(err)

	// google storage
	err = CreateStorageCNAME(ctx, domainName)
	assert.Nil(err)
	exist, err = IsCNAMEExists(ctx, domainName)
	assert.Nil(err)
	assert.True(exist)
	err = DeleteCNAME(ctx, domainName)
	assert.Nil(err)

	TestModeAlwaySuccess()
	defer TestModeBackNormal()

	err = CreateCNAME(ctx, domainName, "ghs.googlehosted.com", false)
	assert.Nil(err)

	exist, err = IsCNAMEExists(ctx, domainName)
	assert.Nil(err)
	assert.True(exist)

	err = DeleteCNAME(ctx, domainName)
	assert.Nil(err)

}

func TestTxtRecord(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	subDomain := "mock-libsrv"
	domainName := subDomain + ".piyuo.com"
	txt := "hi"
	//remove sample record
	RemoveTXT(ctx, domainName)

	exist, err := IsTXTExists(ctx, domainName)
	assert.Nil(err)
	assert.False(exist)

	err = CreateTXT(ctx, domainName, txt)
	assert.Nil(err)

	exist, err = IsTXTExists(ctx, domainName)
	assert.Nil(err)
	assert.True(exist)

	// add txt record that already exist should not error
	err = CreateTXT(ctx, domainName, txt)
	assert.Nil(err)

	err = RemoveTXT(ctx, domainName)
	assert.Nil(err)

	exist, err = IsTXTExists(ctx, domainName)
	assert.Nil(err)
	assert.False(exist)

	// remove txt record second time should not error
	err = RemoveTXT(ctx, domainName)
	assert.Nil(err)

	TestModeAlwaySuccess()
	defer TestModeBackNormal()

	err = CreateTXT(ctx, domainName, txt)
	assert.Nil(err)
	err = RemoveTXT(ctx, domainName)
	assert.Nil(err)
	exist, err = IsTXTExists(ctx, domainName)
	assert.Nil(err)
	assert.True(exist)
}
