package mail

import (
	"context"
	"fmt"

	key "github.com/piyuo/libsrv/key"

	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendgridMail using SendGrid to implement mail
//
type SendgridMail struct {
	BaseMail
}

// newSendgridMail return Mail instance
//
//	mail := newSendgridMail(template)
//
func newSendgridMail(t *template) (Mail, error) {
	mail := &SendgridMail{
		BaseMail: BaseMail{
			Subject:     t.subject,
			Text:        t.text,
			HTML:        t.html,
			FromName:    t.fromName,
			FromAddress: t.fromAddress,
		},
	}
	return mail, nil
}

//Send using sendgrid to send email
//
//	err := mail.Send(ctx)
//
func (c *SendgridMail) Send(ctx context.Context) error {
	sendgridKey, err := key.Text("sendgrid.key")
	if err != nil {
		return err
	}

	m := mail.NewV3Mail()

	from := mail.NewEmail(c.BaseMail.FromName, c.BaseMail.FromAddress)
	m.SetFrom(from)

	textContent := mail.NewContent("text/plain", c.Text)
	m.AddContent(textContent)

	htmlContent := mail.NewContent("text/html", c.HTML)
	m.AddContent(htmlContent)

	personalization := mail.NewPersonalization()
	for _, email := range c.To {
		to := mail.NewEmail(email.Name, email.Address)
		personalization.AddTos(to)
	}
	personalization.Subject = c.Subject
	m.AddPersonalizations(personalization)

	client := sendgrid.NewSendClient(sendgridKey)
	response, err := client.Send(m)
	if err != nil {
		return errors.Wrap(err, "failed to send mail: "+c.Subject)
	}
	// sendgrid status code 2XX is successful send
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return errors.New(fmt.Sprintf("failed to send mail, response=%v, message=%v", response.StatusCode, response.Body))
	}
	//fmt.Print(response)
	return nil
}
