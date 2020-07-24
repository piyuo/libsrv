package data

import (
	"context"
	"encoding/json"
	"time"

	"github.com/piyuo/libsrv/gcp"
	"github.com/piyuo/libsrv/log"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"google.golang.org/api/siteverification/v1"
)

const here = "siteverify"

//VerifiedSite return by SiteVerify.List, show which domain has been verified
//
type VerifiedSite struct {
	ID         string
	DomainName string
}

// SiteVerify verify web site ownership
// must enable site verification api in google cloud project
//
type SiteVerify interface {

	// GetToken return site verify token
	//
	//	ctx := context.Background()
	//	siteverify, err := NewSiteVerify(ctx)
	//	domainName := "mock-site-verify.piyuo.com"
	//	token, err := siteverify.GetToken(ctx, domainName)
	//
	GetToken(ctx context.Context, domainName string) (string, error)

	// Verify return return true if pass verification
	//
	//	ctx := context.Background()
	//	siteverify, err := NewSiteVerify(ctx)
	//	domainName := "mock-site-verify.piyuo.com"
	//	result, err := siteverify.Verify(ctx, domainName)
	//	So(err, ShouldBeNil)
	//	So(result, ShouldBeTrue)
	//
	Verify(ctx context.Context, domainName string) (bool, error)

	// UnVerify return site verify token
	//
	//	ctx := context.Background()
	//	siteverify, err := NewSiteVerify(ctx)
	//	domainName := "mock-site-verify.piyuo.com"
	//	token, err := siteverify.GetToken(ctx, domainName)
	//
	List(ctx context.Context) ([]*VerifiedSite, error)

	// *(this api is not stable,delete domain always error)Delete domain by siteID, you can use List() to get verified site with id
	//
	//	ctx := context.Background()
	//	siteverify, err := NewSiteVerify(ctx)
	//	domainName := "mock-site-verify.piyuo.com"
	//	token, err := siteverify.GetToken(ctx, domainName)
	//
	//Delete(ctx context.Context, siteID string) error
}

// SiteVerifyImpl is cloudflare implementation
//
type SiteVerifyImpl struct {
	SiteVerify
	client       *siteverification.Service
	serviceEmail string
}

// NewSiteVerify create SiteVerify
//
//	ctx := context.Background()
//	storage, err := NewCloudstorage(ctx)
//
func NewSiteVerify(ctx context.Context) (SiteVerify, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	cred, err := gcp.GlobalCredential(ctx)
	if err != nil {
		return nil, err
	}

	client, err := siteverification.NewService(ctx, option.WithCredentials(cred))
	if err != nil {
		return nil, err
	}

	var j map[string]interface{}
	err = json.Unmarshal(cred.JSON, &j)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get credential json")
	}

	siteVerifyImpl := &SiteVerifyImpl{
		client:       client,
		serviceEmail: j["client_email"].(string),
	}

	return siteVerifyImpl, nil
}

// GetToken return site verify token
//
//	ctx := context.Background()
//	siteverify, err := NewSiteVerify(ctx)
//	domainName := "mock-site-verify.piyuo.com"
//	token, err := siteverify.GetToken(ctx, domainName)
//
func (impl *SiteVerifyImpl) GetToken(ctx context.Context, domainName string) (string, error) {
	webService := siteverification.NewWebResourceService(impl.client)
	request := &siteverification.SiteVerificationWebResourceGettokenRequest{
		VerificationMethod: "DNS_TXT",
		Site: &siteverification.SiteVerificationWebResourceGettokenRequestSite{
			Identifier: domainName, //"mock-verify-site.master.piyuo.com"
			Type:       "INET_DOMAIN",
		},
	}

	call := webService.GetToken(request)
	ctx, cancel := context.WithTimeout(ctx, time.Second*12)
	defer cancel()
	call = call.Context(ctx)
	response, err := call.Do()
	if err != nil {
		return "", err
	}
	log.Info(ctx, here, "site verify token created: "+response.Token)
	return response.Token, nil
}

// Verify return return true if pass verification
//
//	ctx := context.Background()
//	siteverify, err := NewSiteVerify(ctx)
//	domainName := "mock-site-verify.piyuo.com"
//	result, err := siteverify.Verify(ctx, domainName)
//	So(err, ShouldBeNil)
//	So(result, ShouldBeTrue)
//
func (impl *SiteVerifyImpl) Verify(ctx context.Context, domainName string) (bool, error) {
	webService := siteverification.NewWebResourceService(impl.client)
	request := &siteverification.SiteVerificationWebResourceResource{
		Site: &siteverification.SiteVerificationWebResourceResourceSite{
			Identifier: domainName,
			Type:       "INET_DOMAIN",
		},
	}

	call := webService.Insert("DNS_TXT", request)
	ctx, cancel := context.WithTimeout(ctx, time.Second*12)
	defer cancel()
	call = call.Context(ctx)
	response, err := call.Do()
	if err != nil {
		return false, err
	}

	for _, v := range response.Owners {
		if v == impl.serviceEmail {
			return true, nil
		}
	}

	return false, nil
}

// List verified site
//
//	ctx := context.Background()
//	siteverify, err := NewSiteVerify(ctx)
//	domainName := "mock-site-verify.piyuo.com"
//	result, err := siteverify.Verify(ctx, domainName)
//	So(err, ShouldBeNil)
//	So(result, ShouldBeTrue)
//
func (impl *SiteVerifyImpl) List(ctx context.Context) ([]*VerifiedSite, error) {
	webService := siteverification.NewWebResourceService(impl.client)
	call := webService.List()
	ctx, cancel := context.WithTimeout(ctx, time.Second*12)
	defer cancel()
	call = call.Context(ctx)
	response, err := call.Do()
	if err != nil {
		return []*VerifiedSite{}, err
	}

	result := make([]*VerifiedSite, len(response.Items))
	for i, v := range response.Items {
		verifiedSite := &VerifiedSite{
			ID:         v.Id,
			DomainName: v.Site.Identifier,
		}
		result[i] = verifiedSite

	}
	return result, nil
}

// *(this api is not stable,delete domain always error)Delete domain by siteID, you can use List() to get verified site with id
//
//	ctx := context.Background()
//	siteverify, err := NewSiteVerify(ctx)
//	domainName := "mock-site-verify.piyuo.com"
//	result, err := siteverify.Verify(ctx, domainName)
//	So(err, ShouldBeNil)
//	So(result, ShouldBeTrue)
//
/*
func (impl *SiteVerifyImpl) Delete(ctx context.Context, siteID string) error {
	webService := siteverification.NewWebResourceService(impl.client)
	call := webService.Delete(siteID)
	ctx, cancel := context.WithTimeout(ctx, time.Second*12)
	defer cancel()
	call = call.Context(ctx)
	err := call.Do()
	if err != nil {
		return err
	}

	return nil
}
*/
