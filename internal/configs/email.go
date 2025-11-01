package configs

type Email struct {
	FromEmail              string `env:"FROM_EMAIL"`
	FromEmailPwd           string `env:"FROM_EMAIL_PASSWORD"`
	FromEmailSMTP          string `env:"FROM_EMAIL_SMTP"`
	SMTPAddr               string `env:"SMTP_ADDR"`
}
