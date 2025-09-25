package router

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers"
	"github.com/labstack/echo/v4"
)

type Router struct {
	echo     *echo.Echo
	handlers *handlers.Handlers
}

func NewRouter(echo *echo.Echo, handlers *handlers.Handlers) *Router {
	return &Router{
		echo:     echo,
		handlers: handlers,
	}
}

func (r *Router) RegisterAPIRoutes() {
	// Health Check
	r.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// v1Public := r.echo.Group("/api/vi")
}
