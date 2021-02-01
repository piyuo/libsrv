package siteverify

import (
	"context"
	"testing"

	cloudflare "github.com/piyuo/libsrv/src/cloudflare"
	"github.com/stretchr/testify/assert"
)

func TestNewSiteVerify(t *testing.T) {
	assert := assert.New(t)
	storage, err := NewSiteVerify(context.Background())
	assert.Nil(err)
	assert.NotNil(storage)
}

func TestVerification(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	siteverify, err := NewSiteVerify(ctx)
	cflare, err := cloudflare.NewCloudflare(ctx)
	domainName := "mock-site-verify.piyuo.com"

	//clean before test
	cflare.RemoveTxtRecord(ctx, domainName)
	defer cflare.RemoveTxtRecord(ctx, domainName)

	token, err := siteverify.GetToken(ctx, domainName)
	assert.Nil(err)
	assert.Greater(len(token), 0)

	//token and token2 should be the same
	token2, err := siteverify.GetToken(ctx, domainName)
	assert.Nil(err)
	assert.Equal(token, token2)

	exist, err := cflare.IsTxtRecordExist(ctx, domainName)
	assert.Nil(err)
	assert.False(exist)
	err = cflare.AddTxtRecord(ctx, domainName, token)
	assert.Nil(err)

	// cause update dns record need time to populate. unmark these test if you want test it manually
	//result, _ := siteverify.Verify(ctx, domainName)
	//assert.Nil(err)
	//assert.True(result)

}
