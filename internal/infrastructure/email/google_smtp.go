package email

import (
	"net/smtp"
	"strings"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/rs/zerolog/log"
)

// type SMTPGoogle interface {
// 	SendEmail(to []string, subject string, body string) error
// }

type SMTPGoogle struct {
	cfg  *configs.Config
	auth smtp.Auth
}

func NewSMTPGoogle(cfg *configs.Config) *SMTPGoogle {
	auth := smtp.PlainAuth("", cfg.Email.FromEmail, cfg.Email.FromEmailPwd, cfg.Email.FromEmailSMTP)
	return &SMTPGoogle{
		cfg:  cfg,
		auth: auth,
	}
}

// SendEmail implements SMTP.
func (s *SMTPGoogle) SendEmail(to []string, subject string, body string) error {
	message := buildMessage(s.cfg.Email.FromEmail, to, subject, body)

	log.Info().Msgf("Send a email to %s", strings.Join(to, ","))
	return smtp.SendMail(
		s.cfg.Email.SMTPAddr,
		s.auth,
		s.cfg.Email.FromEmail,
		to,
		[]byte(message),
	)
}
