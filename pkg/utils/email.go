package utils

import (
	"fmt"
	"strings"
)

// EmailJob represents a single email job (legacy - kept for backward compatibility)
type EmailJob struct {
	Email   string
	Subject string
	Body    string
}

// NewsletterJob represents a batch newsletter job
type NewsletterJob struct {
	UUID      string
	BatchID   int
	Recipient string
	Subject   string
	Body      string
	IsHTML    bool
}

// SubscriberEmail used to retrieve email and uuid of user from the database
type SubscriberEmail struct {
	Uuid  string `json:"uuid,omitempty"`
	Email string `json:"email"`
}

// EmailMessage represents an email message
type EmailMessage struct {
	To      []string
	CC      []string
	BCC     []string
	Subject string
	Body    string
	IsHTML  bool
}

// BuildEmailBody constructs the email body with proper headers
func BuildEmailBody(from string, message *EmailMessage) string {
	var body strings.Builder

	// Headers
	body.WriteString(fmt.Sprintf("From: %s\r\n", from))

	if len(message.To) > 0 {
		body.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(message.To, ", ")))
	}

	if len(message.CC) > 0 {
		body.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(message.CC, ", ")))
	}

	// Note: BCC recipients are not included in headers (that's the point of BCC)

	body.WriteString(fmt.Sprintf("Subject: %s\r\n", message.Subject))

	// Content type
	if message.IsHTML {
		body.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		body.WriteString("MIME-Version: 1.0\r\n")
	} else {
		body.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}

	body.WriteString("\r\n") // Empty line between headers and body
	body.WriteString(message.Body)

	return body.String()
}

// GetAllRecipients combines To, CC, and BCC recipients for SMTP envelope
func GetAllRecipients(message *EmailMessage) []string {
	var recipients []string
	recipients = append(recipients, message.To...)
	recipients = append(recipients, message.CC...)
	recipients = append(recipients, message.BCC...)
	return recipients
}
