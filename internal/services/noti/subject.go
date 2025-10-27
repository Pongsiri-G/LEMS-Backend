package noti

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/events"
)

type NotificationSubject interface {
	Register(observer Observer)
	Deregister(observer Observer)
	Notify(event events.Event)	
}

type notificationSubject struct {
	observers []Observer
}

func NewNotificationSubject() NotificationSubject {
	return &notificationSubject{
		observers: make([]Observer, 0),
	}
}

func (s *notificationSubject) Register(observer Observer) {
	s.observers = append(s.observers, observer)
}

func (s *notificationSubject) Deregister(observer Observer) {
	for i, obs := range s.observers {
		if obs == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func (s *notificationSubject) Notify(event events.Event) {
	for _, obs := range s.observers {
		obs.Update(event)
	}
}
