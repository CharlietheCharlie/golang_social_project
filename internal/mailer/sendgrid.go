package mailer

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (s *SendGridMailer) Send(templateFile, username, email string, data any, isSandBox bool) error {
	from := mail.NewEmail(FromName, s.fromEmail)
	to := mail.NewEmail(username, email)

	// template parsing and building
	template, err := parseTemplate(templateFile, data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, template.Subject, to, "", template.Body)

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandBox,
		},
	})

	return sendWithRetry(maxRetries, email, "SendGrid", func() error {
		response, err := s.client.Send(message)
		if err != nil {
			// transport error / network error
			return err
		}

		if response.StatusCode >= 400 {
			// server error
			return err
		}
		return nil
	})
}
