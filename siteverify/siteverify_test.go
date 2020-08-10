package siteverify

import (
	"context"
	"testing"
	"time"

	cloudflare "github.com/piyuo/libsrv/cloudflare"

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
		cflare, err := cloudflare.NewCloudflare(ctx)
		domainName := "mock-site-verify.piyuo.com"

		//clean before test
		cflare.RemoveTxtRecord(ctx, domainName)

		token, err := siteverify.GetToken(ctx, domainName)
		So(err, ShouldBeNil)
		So(len(token), ShouldBeGreaterThan, 0)

		//token and token2 should be the same
		token2, err := siteverify.GetToken(ctx, domainName)
		So(err, ShouldBeNil)
		So(token, ShouldEqual, token2)

		exist, err := cflare.IsTxtRecordExist(ctx, domainName)
		So(err, ShouldBeNil)
		if !exist {
			err = cflare.AddTxtRecord(ctx, domainName, token)
			So(err, ShouldBeNil)
		}

		// cause update dns record need time to populate. unmark these test if you want test it manually
		//result, err = siteverify.Verify(ctx, domainName)
		//So(err, ShouldBeNil)
		//So(result, ShouldBeTrue)
		time.Sleep(60 * time.Second)
		result, _ := siteverify.Verify(ctx, domainName)
		So(err, ShouldBeNil)
		So(result, ShouldBeTrue)

		//clean after test
		cflare.RemoveTxtRecord(ctx, domainName)
	})

}
