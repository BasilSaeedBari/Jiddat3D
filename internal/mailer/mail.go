package mailer

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendContactNotification sends an email when a new contact form is submitted.
func SendContactNotification(name, method, subject, message string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")
	adminEmail := os.Getenv("ADMIN_EMAIL")

	if host == "" || adminEmail == "" {
		// Silently skip if SMTP is not configured
		return nil
	}

	auth := smtp.PlainAuth("", user, pass, host)

	body := fmt.Sprintf("New contact submission from %s\nContact Method: %s\nSubject: %s\nMessage: %s\n", name, method, subject, message)
	msg := []byte("To: " + adminEmail + "\r\n" +
		"Subject: New Contact Submission: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := fmt.Sprintf("%s:%s", host, port)
	return smtp.SendMail(addr, auth, from, []string{adminEmail}, msg)
}
