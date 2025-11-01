package noti

import (
	"fmt"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/events"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/email"
	"github.com/rs/zerolog/log"
)

type emailObserver struct{
	client email.SMTP
}

type EmailObserver interface {
	Update(event events.Event)
}

func NewEmailObserver(client email.SMTP) EmailObserver {
	return &emailObserver{
		client: client,
	}
}

func (e *emailObserver) Update(event events.Event) {
	payload, ok := event.Payload.(map[string]interface{})
	if !ok {
		return
	}

	userEmail, _ := payload["email"].(string)
	userName, _ := payload["userId"].(string)
	message, _ := payload["message"].(string)

	if userEmail == "" {
		return // skip if no email
	}

	subject := fmt.Sprintf("📢 %s Notification", event.Type)
	body := fmt.Sprintf("Hello %s,\n\n%s\n\n--\nLab Equipment System", userName, message)

	if err := e.client.SendEmail([]string{userEmail}, subject, body); err != nil {
		log.Print("❌ Failed to send email:", err)
		return
	}

	log.Print("📧 Email sent to", userEmail)
}
