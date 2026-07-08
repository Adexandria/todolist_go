package services

type IEmailService interface {
	SendEmail(msg MailgunMessage) (*MailgunResponse, error)
	GenerateHTMLTemplate(templateName string, data map[string]any) (string, error)
	GenerateMailgunTemplate(from string, to string, subject string, html string) MailgunMessage
	SendEmailAsync(msg MailgunMessage)
}
