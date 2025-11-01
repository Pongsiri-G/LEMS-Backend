package log

import (
	logSvc "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/log"
	"github.com/labstack/echo/v4"
)

type LogHandler interface {
	GetAllLogs(c echo.Context) error
}

type handler struct {
	service logSvc.Service
}

func NewLogHandler(service logSvc.Service) LogHandler {
	return &handler{service: service}
}

// GetAllLogs implements Handler.
func (h *handler) GetAllLogs(c echo.Context) error {
	logs, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(500, echo.Map{
			"message": "internal server error",
		})
	}
	return c.JSON(200, logs)
}
