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
	batchSize  = 50 // BCC batch size to avoid SMTP limits
)

type EmailService interface {
	SendEmail(message utils.EmailMessage) error
	SendNewsletter(subscribers []string, subject, body string, isHTML bool) error
}

type SMTPEmailService struct {
	username string
	password string
	smtpHost string
	smtpPort int
}

func NewSMTPEmailService(
	host string,
	port int,
	username string,
	password string,
) EmailService {
	return &SMTPEmailService{
		username: username,
		password: password,
		smtpHost: host,
		smtpPort: port,
	}
}

// SendEmail sends an email using SMTP
func (s *SMTPEmailService) SendEmail(message utils.EmailMessage) error {
	// Create SMTP auth
	auth := smtp.PlainAuth("", s.username, s.password, s.smtpHost)

	// Build the email headers and body
	emailBody := utils.BuildEmailBody(s.username, message)

	recipients := utils.GetAllRecipients(message)

	// Send the email
	addr := s.smtpHost + ":" + strconv.Itoa(s.smtpPort)
	err := smtp.SendMail(addr, auth, s.username, recipients, []byte(emailBody))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendNewsletter sends newsletter to multiple subscribers using BCC batches
func (s *SMTPEmailService) SendNewsletter(subscribers []string, subject, body string, isHTML bool) error {
	// Create batches for BCC sending
	batches := s.createBatches(subscribers, batchSize)

	// Create a job channel
	jobChan := make(chan utils.NewsletterJob, len(batches))

	// WaitGroup to wait for all workers to complete
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go s.newsletterWorker(jobChan, &wg)
	}

	// Send batch jobs to workers
	for i, batch := range batches {
		jobChan <- utils.NewsletterJob{
			BatchID:     i + 1,
			Subscribers: batch,
			Subject:     subject,
			Body:        body,
			IsHTML:      isHTML,
		}
	}

	close(jobChan)
	wg.Wait() // Wait for all workers to finish

	fmt.Printf("Newsletter sent successfully to %d subscribers in %d batches!\n", len(subscribers), len(batches))
	return nil
}

// createBatches splits subscribers into smaller batches for BCC sending
func (s *SMTPEmailService) createBatches(subscribers []string, size int) [][]string {
	var batches [][]string

	for i := 0; i < len(subscribers); i += size {
		end := i + size
		if end > len(subscribers) {
			end = len(subscribers)
		}
		batches = append(batches, subscribers[i:end])
	}

	return batches
}

// newsletterWorker processes newsletter batch jobs
func (s *SMTPEmailService) newsletterWorker(jobs <-chan utils.NewsletterJob, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		message := utils.EmailMessage{
			To:      []string{s.username}, // Send to yourself as primary recipient
			BCC:     job.Subscribers,      // All subscribers in BCC
			Subject: job.Subject,
			Body:    job.Body,
			IsHTML:  job.IsHTML,
		}

		err := s.SendEmail(message)
		if err != nil {
			fmt.Printf("Error sending batch %d: %v\n", job.BatchID, err)
			continue
		}

		fmt.Printf("Successfully sent batch %d to %d recipients\n", job.BatchID, len(job.Subscribers))

		// Rate limiting - delay between batches
		time.Sleep(200 * time.Millisecond)
	}
}
