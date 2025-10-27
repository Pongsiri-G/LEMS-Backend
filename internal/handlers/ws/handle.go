package ws

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/ws"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils/contextutil"
	"github.com/labstack/echo/v4"
)

type WsHandler interface {
	Run(e echo.Context) error
}

type wsHandler struct {
	hub *ws.Hub
}

func NewWsHandler(hub *ws.Hub) WsHandler {
	return &wsHandler{
		hub: hub,
	}
}

// WebSocket implements WsHandler.
func (w *wsHandler) Run(e echo.Context) error {
	user, err := contextutil.GetUserFromContext(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}

	return ws.ServeWS(w.hub, e.Response(), e.Request(), user.ID)
}
