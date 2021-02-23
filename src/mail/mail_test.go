package mail

import (
	"context"
	"net/http"
	"testing"

	"github.com/piyuo/libsrv/src/env"
	"github.com/stretchr/testify/assert"
)

func TestMail(t *testing.T) {
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", "en_US")
	ctx := context.WithValue(context.Background(), env.KeyContextRequest, req)

	mail, err := NewMail(ctx, "mock-mail")
	assert.Nil(err)
	assert.NotEmpty(mail.GetSubject())
	backupSubject := mail.GetSubject()

	mail.SetSubject("ok")
	assert.Equal("ok", mail.GetSubject())
	mail.ReplaceSubject("ok", "1")
	assert.Equal("1", mail.GetSubject())

	assert.NotEmpty(mail.GetText())
	mail.SetText("ok")
	assert.Equal("ok", mail.GetText())
	mail.ReplaceText("ok", "1")
	assert.Equal("1", mail.GetText())

	assert.NotEmpty(mail.GetHTML())
	mail.SetHTML("ok")
	assert.Equal("ok", mail.GetHTML())
	mail.ReplaceHTML("ok", "1")
	assert.Equal("1", mail.GetHTML())

	mail.SetHTML("1")
	mail.SetText("1")
	mail.ReplaceContent("1", "2")
	assert.Equal("2", mail.GetHTML())
	assert.Equal("2", mail.GetText())

	name, address := mail.GetFrom()
	assert.NotEmpty(name)
	assert.NotEmpty(address)
	mail.SetFrom("1", "2")
	name, address = mail.GetFrom()
	assert.Equal("1", name)
	assert.Equal("2", address)

	assert.Nil(mail.GetTo())
	mail.AddTo("name", "address")
	assert.NotNil(mail.GetTo())
	assert.Equal(1, len(mail.GetTo()))
	mail.AddTo("name1", "address1")
	assert.Equal(2, len(mail.GetTo()))
	mail.ResetTo()
	assert.Nil(mail.GetTo())

	//should from cache
	mail, err = NewMail(ctx, "mock-mail")
	assert.Nil(err)
	assert.NotNil(mail)
	assert.Equal(backupSubject, mail.GetSubject())

	mail, err = NewMail(ctx, "notExist")
	assert.NotNil(err)
	assert.Nil(mail)
}

func TestSendMail(t *testing.T) {
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", "en_US")
	ctx := context.WithValue(context.Background(), env.KeyContextRequest, req)
	mail, err := NewMail(ctx, "mock-mail")
	assert.Nil(err)
	mail.AddTo("p", "a@b.c")
	mail.ReplaceText("%1", "1234")
	mail.ReplaceHTML("%1", "1234")
	err = mail.Send(ctx)
	assert.Nil(err)
}

func TestMockSendMail(t *testing.T) {
	assert := assert.New(t)
	TestModeOutputMail = nil
	testMode = true
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", "en_US")
	ctx := context.WithValue(context.Background(), env.KeyContextRequest, req)
	mail, err := NewMail(ctx, "mock-mail")
	assert.Nil(err)
	mail.AddTo("p", "a@b.c")
	mail.ReplaceText("%1", "1234")
	err = mail.Send(ctx)
	assert.Nil(err)
	assert.NotNil(TestModeOutputMail)
	assert.Equal("p", TestModeOutputMail.GetTo()[0].Name)
	assert.Equal("a@b.c", TestModeOutputMail.GetTo()[0].Address)
}
