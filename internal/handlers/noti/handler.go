package noti

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/noti"
	"github.com/labstack/echo/v4"
)

type NotiHandler interface {
	NotifyWebApp(c echo.Context) error
}

type notiHandler struct {
	webNoti noti.WebAppObserver
}

func NewBorrowQueueHandler(webNoti noti.WebAppObserver) NotiHandler {
	return &notiHandler{
		webNoti: webNoti,
	}
}

// NotifyWebApp implements NotiHandler.
func (n *notiHandler) NotifyWebApp(c echo.Context) error {
	panic("not implement")
}
