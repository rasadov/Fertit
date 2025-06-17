package services

import (
	"fmt"
	"github.com/rasadov/MailManagerApp/pkg/utils"
	"net/smtp"
	"strconv"
	"sync"
	"time"
)

const (
	numWorkers = 6
)

type EmailService interface {
	SendEmail(message utils.EmailMessage) error
	SendNewsletter(subscribers []utils.SubscriberEmail, subject, body string, isHTML bool) error
}

type SMTPEmailService struct {
	username string
	password string
	smtpHost string
	smtpPort int
	baseURL  string // Your app's base URL for unsubscribe links
}

func NewSMTPEmailService(
	host string,
	port int,
	username string,
	password string,
	baseURL string,
) EmailService {
	return &SMTPEmailService{
		username: username,
		password: password,
		smtpHost: host,
		smtpPort: port,
		baseURL:  baseURL,
	}
}

// SendEmail sends an email using SMTP
func (s *SMTPEmailService) SendEmail(message utils.EmailMessage) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.smtpHost)
	emailBody := utils.BuildEmailBody(s.username, message)
	recipients := utils.GetAllRecipients(message)

	addr := s.smtpHost + ":" + strconv.Itoa(s.smtpPort)
	err := smtp.SendMail(addr, auth, s.username, recipients, []byte(emailBody))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendNewsletter sends newsletter to multiple subscribers with personalized unsubscribe links
func (s *SMTPEmailService) SendNewsletter(subscribers []utils.SubscriberEmail, subject, body string, isHTML bool) error {
	subscribersLength := len(subscribers)

	if subscribersLength == 0 {
		return fmt.Errorf("no subscribers provided")
	}

	jobChan := make(chan utils.NewsletterJob, subscribersLength)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go s.newsletterWorker(jobChan, &wg)
	}

	// Send individual email jobs with UUID for unsubscribe links
	for i, subscriber := range subscribers {
		// Add unsubscribe link to email body
		personalizedBody := s.addUnsubscribeLink(body, subscriber.Uuid, isHTML)

		jobChan <- utils.NewsletterJob{
			Recipient: subscriber.Email,
			Subject:   subject,
			Body:      personalizedBody,
			IsHTML:    isHTML,
			BatchID:   i + 1,
			UUID:      subscriber.Uuid,
		}
	}

	close(jobChan)
	wg.Wait()

	fmt.Printf("Newsletter sent successfully to %d subscribers!\n", subscribersLength)
	return nil
}

// addUnsubscribeLink adds unsubscribe link to email body
func (s *SMTPEmailService) addUnsubscribeLink(body, uuid string, isHTML bool) string {
	unsubscribeURL := fmt.Sprintf("%s/unsubscribe?token=%s", s.baseURL, uuid)

	if isHTML {
		unsubscribeLink := fmt.Sprintf(`<p><a href="%s">Unsubscribe</a> | <a href="%s">Manage Preferences</a></p>`, unsubscribeURL, unsubscribeURL)
		return body + "\n\n" + unsubscribeLink
	} else {
		unsubscribeText := fmt.Sprintf("\n\nUnsubscribe: %s\nManage Preferences: %s", unsubscribeURL, unsubscribeURL)
		return body + unsubscribeText
	}
}

// newsletterWorker processes individual newsletter jobs
func (s *SMTPEmailService) newsletterWorker(jobs <-chan utils.NewsletterJob, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		message := utils.EmailMessage{
			To:      []string{job.Recipient},
			Subject: job.Subject,
			Body:    job.Body,
			IsHTML:  job.IsHTML,
		}

		err := s.SendEmail(message)
		if err != nil {
			fmt.Printf("Error sending email to %s (job %d): %v\n", job.Recipient, job.BatchID, err)
			continue
		}

		fmt.Printf("Successfully sent email %d to %s\n", job.BatchID, job.Recipient)
		time.Sleep(200 * time.Millisecond)
	}
}
