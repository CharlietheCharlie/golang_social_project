package mailer

import (
	"errors"

	gomail "gopkg.in/gomail.v2"
)

type mailtrapClient struct {
	fromEmail string
	apiKey    string
}

func NewMailTrapClient(apiKey, fromEmail string) (mailtrapClient, error) {
	if apiKey == "" {
		return mailtrapClient{}, errors.New("mailtrap api key is required")
	}
	return mailtrapClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (m mailtrapClient) Send(templateFile, username, email string, data any, isSandBox bool) error {
	// template parsing and building
	tmpl, err := parseTemplate(templateFile, data)
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", tmpl.Subject)
	message.AddAlternative("text/html", tmpl.Body)

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", m.apiKey)

	return sendWithRetry(maxRetries, email, "mailtrap", func() error {
		if err := dialer.DialAndSend(message); err != nil {
			return err
		}
		return nil
	})
}
