package mail

import (
	"context"
	"fmt"
	"log"
	"os"
	app "github.com/piyuo/go-libsrv/app"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

//SendGrid using sendgrid to send email
//
//	mail.Send('from','to','subject','text')
func sgMailOpen(ctx context.Context) {
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}


//SendGrid using sendgrid to send email
//
//	mail.Send('from','to','subject','text')
func SendGrid(ctx context.Context, fromName, fromEmail, toName, toEmail, subject, text, html string) {
	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)
	message := mail.NewSingleEmail(from, subject, to, text, html)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
