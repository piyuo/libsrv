package mail

import (
	"context"

	"github.com/pkg/errors"

	"github.com/piyuo/libsrv/src/i18n"
)

// testMode set to true will put log in test mode. it will print log but not write to database
//
var testMode = false

// EnableTestMode set to true will let every function run success
//
func EnableTestMode(enabled bool) {
	testMode = enabled
}

// TestModeOutputMail is mail sent in test mode
//
var TestModeOutputMail Mail

type template struct {
	subject string

	text string

	html string

	fromName string

	fromAddress string
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
	//	mail.SetFrom("service","service@piyuo.com")
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

// getTemplate get mail template
//
//	template, err := getTemplate("verify", "en_US")
//
func getTemplate(ctx context.Context, name string) (*template, error) {
	json, err := i18n.Resource(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get i18n resource: "+name)
	}

	template := &template{
		subject:     json["subject"].(string),
		text:        json["text"].(string),
		html:        json["html"].(string),
		fromName:    json["fromName"].(string),
		fromAddress: json["fromAddress"].(string),
	}
	return template, nil
}
