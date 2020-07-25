package sms

import (
	"testing"

	"github.com/piyuo/libsrv/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSMS(t *testing.T) {
	Convey("should set/get method", t, func() {
		sms, err := NewSMS("verify", "en-US")
		So(err, ShouldBeNil)
		backupText := sms.GetText()
		So(sms.GetText(), ShouldNotBeEmpty)
		sms.SetText("ok")
		So(sms.GetText(), ShouldEqual, "ok")
		sms.ReplaceText("ok", "1")
		So(sms.GetText(), ShouldEqual, "1")

		//should from cache
		sms, err = NewSMS("verify", "en-US")
		So(err, ShouldBeNil)
		So(sms, ShouldNotBeNil)
		So(sms.GetText(), ShouldEqual, backupText)

	})
}

func TestSMSError(t *testing.T) {
	Convey("should set/get method", t, func() {
		//test not exist template
		sms, err := NewSMS("not exist", "en-US")
		So(err, ShouldNotBeNil)
		So(sms, ShouldBeNil)

		sms, err = NewSMS("verify", "en-US")
		So(err, ShouldBeNil)

		//test canceled ctx
		canceledCtx := util.CanceledCtx()
		err = sms.Send(canceledCtx, "+19999999999")
		So(err, ShouldNotBeNil)
	})
}

func TestSendSMS(t *testing.T) {
	Convey("should send SMS", t, func() {
		sms, err := NewSMS("verify", "en-US")
		So(err, ShouldBeNil)
		sms.ReplaceText("%1", "1234")
		//err = sms.Send(context.Background(), "+19493026176")
		So(err, ShouldBeNil)
	})
}

func TestE164(t *testing.T) {
	Convey("should check number is valid E164 format", t, func() {
		mobile, err := E164("9493017165", "US")
		So(err, ShouldBeNil)
		So(mobile, ShouldEqual, "+19493017165")
		mobile, err = E164("", "")
		So(err, ShouldNotBeNil)
		So(mobile, ShouldBeEmpty)
		mobile, err = E164("94911", "US")
		So(err, ShouldNotBeNil)
		So(mobile, ShouldBeEmpty)
		mobile, err = E164("0987926234", "TW")
		So(err, ShouldBeNil)
		So(mobile, ShouldEqual, "+886987926234")
		mobile, err = E164("9492341654", "TW")
		So(err, ShouldNotBeNil)
		So(mobile, ShouldBeEmpty)
		mobile, err = E164("13916219123", "CN")
		So(err, ShouldBeNil)
		So(mobile, ShouldEqual, "+8613916219123")
		mobile, err = E164("9492341654", "CN")
		So(err, ShouldNotBeNil)
		So(mobile, ShouldBeEmpty)
	})
}
