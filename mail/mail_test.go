package mail

import (
	"context"
	"net/http"
	"testing"

	"github.com/piyuo/libsrv/env"
	"github.com/piyuo/libsrv/google"
	"github.com/stretchr/testify/assert"
)

func TestMail(t *testing.T) {
	t.Parallel()
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

	// from cache
	mail, err = NewMail(ctx, "mock-mail")
	assert.Nil(err)
	assert.NotNil(mail)
	assert.Equal(backupSubject, mail.GetSubject())

	mail, err = NewMail(ctx, "notExist")
	assert.NotNil(err)
	assert.Nil(mail)
}

func TestSendMail(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", "en_US")
	ctx := context.WithValue(context.Background(), env.KeyContextRequest, req)
	mail, err := NewMail(ctx, "mock-mail")
	assert.Nil(err)
	mail.AddTo(google.TestProject, google.TestEmail)
	mail.ReplaceContent("%1%", "1234")
	err = mail.Send(ctx)
	assert.Nil(err)
}

func TestMock(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", "en_US")
	ctx := context.WithValue(context.Background(), env.KeyContextRequest, req)
	mail, err := NewMail(ctx, "mock-mail")
	assert.Nil(err)
	//	mail.AddTo(google.TestProject, google.TestEmail)
	mail.AddTo(google.TestProject, "791088@gmail.com")
	mail.ReplaceContent("%1%", "123456")

	ctx = context.WithValue(context.Background(), MockSuccess, "")
	err = mail.Send(ctx)
	assert.Nil(err)

	ctx = context.WithValue(context.Background(), MockError, "")
	err = mail.Send(ctx)
	assert.NotNil(err)

	ctx = context.WithValue(context.Background(), KeepMail, "")
	LastMail = nil
	mail.Send(ctx)
	assert.NotNil(LastMail)
	assert.Equal(google.TestProject, LastMail.GetTo()[0].Name)
	//	assert.Equal(google.TestEmail, LastMail.GetTo()[0].Address)
}

func TestGetTemplate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept-Language", "en_US")
	ctx := context.WithValue(context.Background(), env.KeyContextRequest, req)

	template, err := getTemplate(ctx, "mock-mail")
	assert.Nil(err)
	assert.NotNil(template)

	// get template content from cache
	template2, err := getTemplate(ctx, "mock-mail")
	assert.Nil(err)
	assert.NotNil(template2)

	assert.Equal(template.HTML, template2.HTML)
}

func TestLocalizedContent(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	j := map[string]interface{}{
		"{1}": "112233",
	}

	content := localizedContent("mock-mail.txt", j)
	assert.NotEmpty(content)
	assert.Contains(content, "112233")
}

func BenchmarkGetTemplate(b *testing.B) {
	ctx := context.Background()
	getTemplate(ctx, "mock-email")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getTemplate(ctx, "mock-email")
	}
}
