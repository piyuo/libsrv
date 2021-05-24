package mail

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/piyuo/libsrv/file"
	"github.com/piyuo/libsrv/i18n"
	"github.com/piyuo/libsrv/mapping"
)

const (
	CacheKey = "m-"
)

// Mock define key test flag
//
type Mock int8

const (
	// MockSuccess let function return nil
	//
	MockSuccess Mock = iota

	// MockError let function error
	//
	MockError

	// KeepMail keep mail in LastMail
	//
	KeepMail
)

// LastMail is mail sent when KeepMail
//
var LastMail Mail

// forceStopSend is true will stop send email
//
var forceStopSend = false

// ForceStopSend set to true will stop send email
//
func ForceStopSend(value bool) {
	forceStopSend = value
}

type Template struct {
	Subject string

	Text string

	HTML string

	FromName string

	FromAddress string
}

// Mail use template to generate mail content and send
//
type Mail interface {

	// GetSubject return mail subject
	//
	//	subject := mail.GetSubject()
	//
	GetSubject() string

	// SetSubject set mail subject
	//
	//	mail.SetSubject("subject")
	//
	SetSubject(subject string) *BaseMail

	// GetText return mail text content
	//
	//	text := mail.GetText()
	//
	GetText() string

	// SetText set mail text content
	//
	//	mail.SetText("text body")
	//
	SetText(text string) *BaseMail

	// GetHTML return mail html content
	//
	//	html := mail.GetHTML()
	//
	GetHTML() string

	// SetHTML set mail html content
	//
	//	mail.SetHTML("html body")
	//
	SetHTML(html string) *BaseMail

	// GetFrom return from email address
	//
	//	name,address := mail.GetFrom()
	//
	GetFrom() (string, string)

	// SetFrom set from email address
	//
	//	mail.SetFrom("service","service@somedomain.com")
	//
	SetFrom(emailName, emailAddress string) *BaseMail

	// ReplaceSubject replace string in mail subject
	//
	//	mail.ReplaceSubject("%1","hello")
	//
	ReplaceSubject(replaceFrom, replaceTo string) *BaseMail

	// ReplaceContent replace string in mail text and html content
	//
	//	mail.ReplaceContent("%1","hello")
	//
	ReplaceContent(replaceFrom, replaceTo string) *BaseMail

	// ReplaceText replace string in mail text content
	//
	//	mail.ReplaceText("%1","hello")
	//
	ReplaceText(replaceFrom, replaceTo string) *BaseMail

	// ReplaceHTML replace string in html comtent
	//
	//	to := mail.GetTo()
	//
	ReplaceHTML(replaceFrom, replaceTo string) *BaseMail

	// GetTo get email to
	//
	//	to := mail.GetTo()
	//
	GetTo() []*Email

	// AddTo add email to
	//
	//	mail.AddTo("user","a@b.c")
	//
	AddTo(emailName, emailAddress string) *BaseMail

	// ResetTo reset to empty
	//
	//	mail.ResetTo()
	//
	ResetTo() *BaseMail

	// Send mail
	//
	//	m, err := mail.NewMail("verify", "en_US;'")
	//	m.AddTo("piyuo", "a@b.c")
	//	m.ReplaceText("%1", "1234")
	//	m.ReplaceHTML("%1", "1234")
	//	err := m.Send(ctx)
	//
	Send(ctx context.Context) error
}

// Email is a single email address
//
type Email struct {

	// Name is email name
	//
	Name string

	// Address is email address
	//
	Address string
}

// NewMail return Mail instance, require template name and locale to find template
//
//	m, err := mail.NewMail("verify", "en_US")
//	m.AddTo("piyuo", "a@b.c")
//	m.ReplaceText("%1", "1234")
//	m.ReplaceHTML("%1", "1234")
//	m.Send(ctx)
//
func NewMail(ctx context.Context, name string) (Mail, error) {
	template, err := getTemplate(ctx, name)
	if err != nil {
		return nil, err
	}
	return newSendgridMail(template)
}

// getTemplate get mail template, template will be cache for 24 hour
//
//	template, err := getTemplate("mock")
//
func getTemplate(ctx context.Context, name string) (*Template, error) {
	// BenchmarkGetTemplate-16    	   45402	     24267 ns/op	    4280 B/op	      61 allocs/op
	jsonContent, err := i18n.JSON(ctx, name, ".json", 24*time.Hour)
	if err != nil {
		return nil, errors.Wrapf(err, "i18n json %v", name)
	}

	// html don't have locale
	htmlContent, err := file.I18nText(name+".html", 24*time.Hour)
	if err != nil {
		return nil, errors.Wrapf(err, "i18n text %v", name)
	}

	// replace htmlContent tag {xxx} from json
	for key, value := range jsonContent {
		if key[0] == '{' {
			htmlContent = strings.ReplaceAll(htmlContent, key, fmt.Sprint(value))
		}
	}

	// don't cache template it will be slower due to digit.encode() 
	template := &Template{
		Subject:     mapping.GetString(jsonContent, "subject", ""),
		Text:        mapping.GetString(jsonContent, "text", ""),
		FromName:    mapping.GetString(jsonContent, "fromName", ""),
		FromAddress: mapping.GetString(jsonContent, "fromAddress", ""),
		HTML:        htmlContent,
	}
	return template, nil
}
