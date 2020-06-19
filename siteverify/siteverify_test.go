package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewSiteVerify(t *testing.T) {
	Convey("should new SiteVerify", t, func() {
		storage, err := NewSiteVerify(context.Background())
		So(err, ShouldBeNil)
		So(storage, ShouldNotBeNil)
	})
}

func TestVerification(t *testing.T) {
	Convey("should verify domain", t, func() {
		ctx := context.Background()
		siteverify, err := NewSiteVerify(ctx)
		domainName := "mock-site-verify.piyuo.com"

		sites, err := siteverify.List(ctx)
		for _, site := range sites {
			siteverify.Delete(ctx, site.ID)
		}

		result, err := siteverify.Verify(ctx, domainName)
		So(err, ShouldNotBeNil)

		token, err := siteverify.GetToken(ctx, domainName)
		So(err, ShouldBeNil)
		So(len(token), ShouldBeGreaterThan, 0)

		//		result, err := siteverify.Verify(ctx, domainName)
		So(err, ShouldBeNil)
		So(result, ShouldBeTrue)

	})

}
