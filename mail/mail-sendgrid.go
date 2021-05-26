package mail

import (
	"context"
	"fmt"

	"github.com/piyuo/libsrv/file"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var sendgridClient *sendgrid.Client

func getSendgridClient() (*sendgrid.Client, error) {
	if sendgridClient != nil {
		return sendgridClient, nil
	}
	sendgridKey, err := file.KeyText("sendgrid.key")
	if err != nil {
		return nil, errors.Wrap(err, "get key")
	}
	sendgridClient = sendgrid.NewSendClient(sendgridKey)
	return sendgridClient, nil
}

// SendgridMail using SendGrid to implement mail
//
type SendgridMail struct {
	BaseMail
}

// newSendgridMail return Mail instance
//
//	mail := newSendgridMail(template)
//
func newSendgridMail(t *Template) (Mail, error) {
	mail := &SendgridMail{
		BaseMail: BaseMail{
			Subject: t.Subject,
			Text:    t.Text,
			HTML:    t.HTML,
			Sender:  t.Sender,
			From:    t.From,
		},
	}
	return mail, nil
}

//Send using sendgrid to send email
//
//	err := mail.Send(ctx)
//
func (c *SendgridMail) Send(ctx context.Context) error {
	if ctx.Value(KeepMail) != nil {
		LastMail = c
	}
	if forceStopSend || ctx.Value(MockSuccess) != nil {
		return nil
	}
	if ctx.Value(MockError) != nil {
		return errors.New("")
	}

	m := mail.NewV3Mail()
	from := mail.NewEmail(c.BaseMail.Sender, c.BaseMail.From)
	m.SetFrom(from)

	textContent := mail.NewContent("text/plain", c.Text)
	m.AddContent(textContent)

	if c.HTML != "" {
		htmlContent := mail.NewContent("text/html", c.HTML)
		m.AddContent(htmlContent)
	}

	personalization := mail.NewPersonalization()
	for _, email := range c.To {
		to := mail.NewEmail(email.Name, email.Address)
		personalization.AddTos(to)
	}
	personalization.Subject = c.Subject
	m.AddPersonalizations(personalization)

	client, err := getSendgridClient()
	if err != nil {
		return errors.Wrapf(err, "get client")
	}

	response, err := client.Send(m)
	if err != nil {
		return errors.Wrapf(err, "sendgrid fail %v", c.Subject)
	}
	// sendgrid status code 2XX is successful send
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return errors.New(fmt.Sprintf("sendgrid error, response=%v, message=%v", response.StatusCode, response.Body))
	}
	return nil
}
