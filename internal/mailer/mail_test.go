package mailer

import "testing"

func TestSanitizeHeader(t *testing.T) {
	got := sanitizeHeader("Hello\r\nBcc: evil@example.com\n")
	want := "HelloBcc: evil@example.com"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
