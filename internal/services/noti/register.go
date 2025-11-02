package noti

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/email"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/ws"
)

func ProvideSubjectWithObservers(subject NotificationSubject, hub *ws.Hub, smtp *email.SMTPGoogle) Subject {
	subject.Register(NewWebAppObserver(hub))
	subject.Register(NewEmailObserver(smtp))
	return subject
}