package services

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type MailgunMessage struct {
	From        string
	To          string
	CC          string
	BCC         string
	Subject     string
	Text        string
	HTML        string
	Template    string
	TemplateVar string
	Tag         string
}

type MailgunResponse struct {
	StatusCode int
	Body       string
}

var templateVars = map[string]string{
	"Welcome": "Templates/welcome.html",
	"Login":   "Templates/login.html",
	"Notify":  "Templates/notification.html",
	"Forgot":  "Templates/forgot_password.html",
}

type EmailService struct {
}

func (e *EmailService) GenerateMailgunTemplate(from string, to string, subject string, html string) MailgunMessage {
	return MailgunMessage{
		From:    from,
		To:      to,
		Subject: subject,
		HTML:    html,
	}
}

func EmailServiceCon() *EmailService {
	return &EmailService{}
}

func (e *EmailService) SendEmail(msg MailgunMessage) (*MailgunResponse, error) {
	domainName := os.Getenv("DOMAIN_NAME")
	apiKey := os.Getenv("API_KEY")
	reqURL := "https://api.mailgun.net/v3/" + domainName + "/messages"

	fields := map[string]string{
		"from":        msg.From,
		"to":          msg.To,
		"cc":          msg.CC,
		"bcc":         msg.BCC,
		"subject":     msg.Subject,
		"text":        msg.Text,
		"html":        msg.HTML,
		"template":    msg.Template,
		"t:variables": msg.TemplateVar,
		"t:text":      "yes",
		"o:dkim":      "yes",
		"o:tracking":  "yes",
		"o:tag":       msg.Tag,
	}

	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)

	for field, value := range fields {
		if value == "" {
			continue
		}
		fw, err := writer.CreateFormField(field)
		if err != nil {
			return nil, fmt.Errorf("creating form field %q: %w", field, err)
		}
		if _, err = io.Copy(fw, strings.NewReader(value)); err != nil {
			return nil, fmt.Errorf("writing form field %q: %w", field, err)
		}
	}
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewReader(data.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.SetBasicAuth("api", apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	return &MailgunResponse{
		StatusCode: res.StatusCode,
		Body:       string(body),
	}, nil
}

func (e *EmailService) SendEmailAsync(msg MailgunMessage) {
	go func() {
		res, err := e.SendEmail(msg)
		if err != nil {
			slog.Error("failed to send email", "error", err, "to", msg.To)
			return
		}
		if res.StatusCode >= 400 {
			slog.Warn("email send failed", "status", res.StatusCode, "body", res.Body, "to", msg.To)
		}
	}()
}

func (e *EmailService) GenerateHTMLTemplate(templateName string, data map[string]any) (string, error) {
	templateLocation := templateVars[templateName]
	content, err := os.ReadFile(templateLocation)
	if err != nil {
		return "", fmt.Errorf("reading template file: %w", err)
	}

	t, err := template.New("template").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

var _ IEmailService = (*EmailService)(nil)
