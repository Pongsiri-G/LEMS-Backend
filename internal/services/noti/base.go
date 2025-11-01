package noti

import "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/events"

type Observer interface {
	Update(event events.Event)
}

type Subject interface {
	Register(observer Observer)
	Deregister(observer Observer)
	Notify(event events.Event)
}
