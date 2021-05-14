package gsite

import (
	"context"
	"testing"

	cloudflare "github.com/piyuo/libsrv/cloudflare"
	"github.com/piyuo/libsrv/google"
	"github.com/stretchr/testify/assert"
)

func TestNewSiteVerify(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	storage, err := NewSiteVerify(context.Background())
	assert.Nil(err)
	assert.NotNil(storage)
}

func TestVerification(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	ctx := context.Background()
	siteverify, err := NewSiteVerify(ctx)
	assert.Nil(err)
	domainName := "mock-site-verify." + google.MyDomain

	//clean before test
	cloudflare.RemoveTXT(ctx, domainName)
	defer cloudflare.RemoveTXT(ctx, domainName)

	token, err := siteverify.GetToken(ctx, domainName)
	assert.Nil(err)
	assert.Greater(len(token), 0)

	//token and token2 should be the same
	token2, err := siteverify.GetToken(ctx, domainName)
	assert.Nil(err)
	assert.Equal(token, token2)

	exist, err := cloudflare.IsTXTExists(ctx, domainName)
	assert.Nil(err)
	assert.False(exist)
	err = cloudflare.CreateTXT(ctx, domainName, token)
	assert.Nil(err)

	// cause update dns record need time to populate. unmark these test if you want test it manually
	//result, _ := siteverify.Verify(ctx, domainName)
	//assert.Nil(err)
	//assert.True(result)

}
