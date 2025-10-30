package ws

import (
	"fmt"
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/events"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/ws"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/noti"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type WsHandler interface {
	Run(e echo.Context) error
	SendNoti(e echo.Context) error
}

type wsHandler struct {
	hub *ws.Hub
	events     noti.Subject
}

func NewWsHandler(hub *ws.Hub, events noti.Subject) WsHandler {
	return &wsHandler{
		hub: hub,
		events: events,
	}
}

// WebSocket implements WsHandler.
func (w *wsHandler) Run(e echo.Context) error {
	userId := e.Request().URL.Query().Get("userId")

	return ws.ServeWS(w.hub, e.Response(), e.Request(), userId)
}



func (w *wsHandler) SendNoti(e echo.Context) error {
	var req requests.BorrowRequest
	if err := e.Bind(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind create tag request")
		return e.JSON(http.StatusBadRequest, nil)
	}

	w.events.Notify(events.Event{
		Type: "Demo Notification",
		Payload: map[string]interface{}{
			"userId":      req.UserID,
			"equipmentId": req.ItemID,
			"message":     fmt.Sprintf("Notify to you (%s)", req.UserID),
		},

	})	
	return e.JSON(http.StatusOK, map[string]string{ "message": "success" })
}
