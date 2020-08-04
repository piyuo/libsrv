package cloudflare

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

func TestDomain(t *testing.T) {
	Convey("should add remove sub domain", t, func() {
		ctx := context.Background()
		cflare, err := NewCloudflare(ctx)
		So(err, ShouldBeNil)
		So(cflare, ShouldNotBeNil)
		subDomain := "mock-libsrv"
		domainName := subDomain + ".piyuo.com"

		//remove sample domain
		cflare.RemoveDomain(ctx, domainName)

		exist, err := cflare.IsDomainExist(ctx, domainName)
		So(err, ShouldBeNil)
		So(exist, ShouldBeFalse)

		err = cflare.AddDomain(ctx, domainName, false)
		So(err, ShouldBeNil)

		exist, err = cflare.IsDomainExist(ctx, domainName)
		So(err, ShouldBeNil)
		So(exist, ShouldBeTrue)

		// add domain that already exist should not error
		err = cflare.AddDomain(ctx, domainName, false)
		So(err, ShouldBeNil)

		err = cflare.RemoveDomain(ctx, domainName)
		So(err, ShouldBeNil)

		exist, err = cflare.IsDomainExist(ctx, domainName)
		So(err, ShouldBeNil)
		So(exist, ShouldBeFalse)

		// remove domain second time should not error
		err = cflare.RemoveDomain(ctx, domainName)
		So(err, ShouldBeNil)

	})
}

func TestTxtRecord(t *testing.T) {
	Convey("should add remove sub domain", t, func() {
		ctx := context.Background()
		cflare, err := NewCloudflare(ctx)
		So(err, ShouldBeNil)
		So(cflare, ShouldNotBeNil)
		subDomain := "mock-libsrv"
		domainName := subDomain + ".piyuo.com"
		txt := "hi"
		//remove sample record
		cflare.RemoveTxtRecord(ctx, domainName)

		exist, err := cflare.IsTxtRecordExist(ctx, domainName)
		So(err, ShouldBeNil)
		So(exist, ShouldBeFalse)

		err = cflare.AddTxtRecord(ctx, domainName, txt)
		So(err, ShouldBeNil)

		exist, err = cflare.IsTxtRecordExist(ctx, domainName)
		So(err, ShouldBeNil)
		So(exist, ShouldBeTrue)

		// add txt record that already exist should not error
		err = cflare.AddTxtRecord(ctx, domainName, txt)
		So(err, ShouldBeNil)

		err = cflare.RemoveTxtRecord(ctx, domainName)
		So(err, ShouldBeNil)

		exist, err = cflare.IsTxtRecordExist(ctx, domainName)
		So(err, ShouldBeNil)
		So(exist, ShouldBeFalse)

		// remove txt record second time should not error
		err = cflare.RemoveTxtRecord(ctx, domainName)
		So(err, ShouldBeNil)

	})
}
