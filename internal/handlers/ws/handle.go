package ws

import (
	"fmt"
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/events"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/ws"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/noti"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/user"
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
	userService user.UserService
}

func NewWsHandler(hub *ws.Hub, events noti.Subject, userService user.UserService) WsHandler {
	return &wsHandler{
		hub: hub,
		events: events,
		userService: userService,
	}
}

// WebSocket implements WsHandler.
func (w *wsHandler) Run(e echo.Context) error {
	userId := e.Request().URL.Query().Get("userId")

	return ws.ServeWS(w.hub, e.Response(), e.Request(), userId)
}



func (w *wsHandler) SendNoti(e echo.Context) error {
	var req struct{
		UserID string	`json:"user_id"`
	}

	if err := e.Bind(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind create tag request")
		return e.JSON(http.StatusBadRequest, nil)
	}

	user, err := w.userService.FindByID(e.Request().Context(), req.UserID)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{ "message" : err.Error()})
	}

	if user == nil {
		return e.JSON(http.StatusNotFound, map[string]string{ "message" : exceptions.ErrUserNotFound.Error()})
	}

	w.events.Notify(events.Event{
		Type: "Demo Notification",
		Payload: map[string]interface{}{
			"userId":      req.UserID,
			"message":     fmt.Sprintf("Notify to you (%s)", req.UserID),
			"email": 	   user.UserEmail,
		},

	})	
	return e.JSON(http.StatusOK, map[string]string{ "message": "success" })
}
