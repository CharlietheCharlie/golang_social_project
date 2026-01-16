package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"time"
)

const (
	FromName            = "social"
	maxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed templates
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandBox bool) error
}

type Template struct {
	Subject string
	Body    string
}

func parseTemplate(templateFile string, data any) (Template, error) {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return Template{}, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return Template{}, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return Template{}, err
	}

	return Template{
		Subject: subject.String(),
		Body:    body.String(),
	}, nil
}

func sendWithRetry(maxRetries int, email, provider string, sendFn func() error) error {
	for i := 0; i < maxRetries; i++ {
		err := sendFn()
		if err == nil {
			log.Printf("[%s] Email sent to %s", provider, email)
			return nil
		}

		log.Printf("[%s] Failed to send email to %s (attempt %d/%d): %v",
			provider, email, i+1, maxRetries, err)

		// 最後一次就別 sleep 了
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second) // exponential-ish backoff
		}
	}

	return fmt.Errorf("failed to send email to %s after %d attempts", email, maxRetries)
}
