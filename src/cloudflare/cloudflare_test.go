package cloudflare

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCloudflare(t *testing.T) {
	assert := assert.New(t)
	cflare, err := NewCloudflare(context.Background())
	assert.Nil(err)
	assert.NotNil(cflare)
}

func TestDomain(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	cflare, err := NewCloudflare(ctx)
	assert.Nil(err)
	assert.NotNil(cflare)
	subDomain := "mock-libsrv"
	domainName := subDomain + ".piyuo.com"

	//remove sample domain
	cflare.RemoveDomain(ctx, domainName)

	exist, err := cflare.IsDomainExist(ctx, domainName)
	assert.Nil(err)
	assert.False(exist)

	err = cflare.AddDomain(ctx, domainName, false)
	assert.Nil(err)

	exist, err = cflare.IsDomainExist(ctx, domainName)
	assert.Nil(err)
	assert.True(exist)

	// add domain that already exist should not error
	err = cflare.AddDomain(ctx, domainName, false)
	assert.Nil(err)

	err = cflare.RemoveDomain(ctx, domainName)
	assert.Nil(err)

	exist, err = cflare.IsDomainExist(ctx, domainName)
	assert.Nil(err)
	assert.False(exist)

	// remove domain second time should not error
	err = cflare.RemoveDomain(ctx, domainName)
	assert.Nil(err)
}

func TestTxtRecord(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	cflare, err := NewCloudflare(ctx)
	assert.Nil(err)
	assert.NotNil(cflare)
	subDomain := "mock-libsrv"
	domainName := subDomain + ".piyuo.com"
	txt := "hi"
	//remove sample record
	cflare.RemoveTxtRecord(ctx, domainName)

	exist, err := cflare.IsTxtRecordExist(ctx, domainName)
	assert.Nil(err)
	assert.False(exist)

	err = cflare.AddTxtRecord(ctx, domainName, txt)
	assert.Nil(err)

	exist, err = cflare.IsTxtRecordExist(ctx, domainName)
	assert.Nil(err)
	assert.True(exist)

	// add txt record that already exist should not error
	err = cflare.AddTxtRecord(ctx, domainName, txt)
	assert.Nil(err)

	err = cflare.RemoveTxtRecord(ctx, domainName)
	assert.Nil(err)

	exist, err = cflare.IsTxtRecordExist(ctx, domainName)
	assert.Nil(err)
	assert.False(exist)

	// remove txt record second time should not error
	err = cflare.RemoveTxtRecord(ctx, domainName)
	assert.Nil(err)
}
