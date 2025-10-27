package noti

import (
	"encoding/json"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/events"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/ws"
	"github.com/rs/zerolog/log"
)

type WebAppObserver interface {
	Update(event events.Event)
}
type webAppObserver struct {
	hub *ws.Hub
}

func NewWebAppObserver(hub *ws.Hub) WebAppObserver {
	return &webAppObserver{
		hub: hub,
	}
}

type pushMsg struct {
	Type    string `json:"type"`
	Message string `json:"massage"`
	UserID  string `json:"userId,omitempty"`
}

func (w *webAppObserver) Update(event events.Event) {
	msg := pushMsg{
		Type: string(event.Type),
		Message: "Notification",
	}

	if m, ok := event.Payload.(map[string]any); ok {
		if uid, ok2 := m["userId"].(string); ok2 {
			msg.UserID = uid
		}
		if text, ok2 := m["message"].(string); ok2 {
			msg.Message = text
		}
	}

	if b, err := json.Marshal(msg); err == nil {
		w.hub.SendToUser(msg.UserID, b)
	}

	log.Info().Msgf("%s success", string(event.Type))
}
