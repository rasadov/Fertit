package services

type EmailService interface{}

func NewSMTPEmailService(
	host string,
	port int,
	username string,
	password string,
) *EmailService {
	return nil
}
