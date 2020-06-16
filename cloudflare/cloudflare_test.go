package data

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewCloudflare(t *testing.T) {
	Convey("should new cloudflare", t, func() {
		cflare, err := NewCloudflare(context.Background())
		So(err, ShouldBeNil)
		So(cflare, ShouldNotBeNil)
	})
}

func TestSubDomain(t *testing.T) {
	Convey("should add remove sub domain", t, func() {
		ctx := context.Background()
		cflare, err := NewCloudflare(ctx)
		So(err, ShouldBeNil)
		So(cflare, ShouldNotBeNil)
		subDomain := "mock-libsrv"
		domainName := subDomain + ".piyuo.com"

		cflare.RemoveDomain(ctx, domainName)

		exist, err := cflare.IsDomainExist(ctx, domainName)
		So(err, ShouldBeNil)
		So(exist, ShouldBeFalse)

		err = cflare.AddDomain(ctx, domainName, false)
		So(err, ShouldBeNil)

		exist, err = cflare.IsDomainExist(ctx, domainName)
		So(err, ShouldBeNil)
		So(exist, ShouldBeTrue)

		err = cflare.RemoveDomain(ctx, domainName)
		So(err, ShouldBeNil)
	})
}
