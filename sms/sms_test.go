package sms

import (
	"testing"

	"github.com/piyuo/libsrv/test"
	"github.com/stretchr/testify/assert"
)

func TestSMS(t *testing.T) {
	assert := assert.New(t)
	sms, err := NewSMS("verify", "en-US")
	assert.Nil(err)
	backupText := sms.GetText()
	assert.NotEmpty(sms.GetText())
	sms.SetText("ok")
	assert.Equal("ok", sms.GetText())
	sms.ReplaceText("ok", "1")
	assert.Equal("1", sms.GetText())

	//should from cache
	sms, err = NewSMS("verify", "en-US")
	assert.Nil(err)
	assert.NotNil(sms)
	assert.Equal(backupText, sms.GetText())
}

func TestSMSError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	//test not exist template
	sms, err := NewSMS("not exist", "en-US")
	assert.NotNil(err)
	assert.Nil(sms)

	sms, err = NewSMS("verify", "en-US")
	assert.Nil(err)

	//test canceled ctx
	canceledCtx := test.CanceledContext()
	err = sms.Send(canceledCtx, "+19999999999")
	assert.NotNil(err)
}

func TestSendSMS(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	sms, err := NewSMS("verify", "en-US")
	assert.Nil(err)
	sms.ReplaceText("%1", "1234")
	//err = sms.Send(context.Background(), "+19493026176")
	assert.Nil(err)
}

func TestE164(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	mobile, err := E164("9493017165", "US")
	assert.Nil(err)
	assert.Equal("+19493017165", mobile)
	mobile, err = E164("", "")
	assert.NotNil(err)
	assert.Empty(mobile)
	mobile, err = E164("94911", "US")
	assert.NotNil(err)
	assert.Empty(mobile)
	mobile, err = E164("0987926234", "TW")
	assert.Nil(err)
	assert.Equal("+886987926234", mobile)
	mobile, err = E164("9492341654", "TW")
	assert.NotNil(err)
	assert.Empty(mobile)
	mobile, err = E164("13916219123", "CN")
	assert.Nil(err)
	assert.Equal("+8613916219123", mobile)
	mobile, err = E164("9492341654", "CN")
	assert.NotNil(err)
	assert.Empty(mobile)
}
