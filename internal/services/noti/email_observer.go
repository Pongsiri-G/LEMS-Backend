package noti

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/events"
)

type emailObserver struct{}

func NewEmailObserver() Observer {
	return &emailObserver{}
}

func (e *emailObserver) Update(event events.Event) {
}
