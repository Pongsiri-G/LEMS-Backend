package server

import (
	"fmt"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers"
	"github.com/labstack/echo/v4"
	// echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type EchoServer struct {
	config *configs.Config
	handlers       *handlers.Handlers
	// authMiddleware middlewares.AuthMiddleware
}

func NewEchoServer(
	config *configs.Config,
	handlers *handlers.Handlers,
	// authMiddleware middlewares.AuthMiddleware,
) *EchoServer {
	return &EchoServer{
		config: config,
		handlers: handlers,	
	}
}

func (s *EchoServer) Start() error {
	e := echo.New()

	// e.Validator = validator.NewValidator()

	// e.HTTPErrorHandler = servererr.EchoHTTPErrorHandler

	// e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
	// 	AllowOrigins:     s.config.CORS.AllowOrigins,
	// 	AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
	// 	AllowCredentials: true,
	// 	AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	// }))

	// router := router.NewRouter(e, s.handlers, s.authMiddleware)

	// router.RegisterAPIRoutes()

	return e.Start(fmt.Sprintf(":%s", s.config.Port))
}
