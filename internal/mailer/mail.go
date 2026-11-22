package mailer

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

func sanitizeHeader(s string) string {
	return strings.NewReplacer("\r", "", "\n", "").Replace(strings.TrimSpace(s))
}

// SendContactNotification sends an email when a new contact form is submitted.
func SendContactNotification(name, method, subject, message string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_SENDER_EMAIL")
	if from == "" {
		from = os.Getenv("SMTP_FROM")
	}
	if from == "" {
		from = user
	}
	adminEmail := sanitizeHeader(os.Getenv("ADMIN_EMAIL"))
	from = sanitizeHeader(from)

	if host == "" || adminEmail == "" || from == "" {
		// Silently skip if SMTP is not configured
		return nil
	}

	auth := smtp.PlainAuth("", user, pass, host)

	safeSubject := sanitizeHeader(subject)
	body := fmt.Sprintf("New contact submission from %s\nContact Method: %s\nSubject: %s\nMessage: %s\n", name, method, safeSubject, message)
	msg := []byte("From: " + from + "\r\n" +
		"To: " + adminEmail + "\r\n" +
		"Subject: New Contact Submission: " + safeSubject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := fmt.Sprintf("%s:%s", host, port)
	return smtp.SendMail(addr, auth, from, []string{adminEmail}, msg)
}
