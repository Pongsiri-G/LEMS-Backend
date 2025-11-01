package email

import (
	"fmt"
	"strings"
)

type SMTP interface {
	SendEmail(to []string, subject, body string) error
}

func buildMessage(from string, to []string, subject, body string) string {
	return fmt.Sprintf(
		"From: System <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s\r\n",
		from,
		strings.Join(to, ", "),
		subject,
		body,
	)
}
