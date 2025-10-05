package router

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/middlewares"
	"github.com/labstack/echo/v4"
)

type Router struct {
	echo           *echo.Echo
	handlers       *handlers.Handlers
	authMiddleware middlewares.AuthMiddleware
}

func NewRouter(echo *echo.Echo, handlers *handlers.Handlers, authMiddleware middlewares.AuthMiddleware) *Router {
	return &Router{
		echo:           echo,
		handlers:       handlers,
		authMiddleware: authMiddleware,
	}
}

func (r *Router) RegisterAPIRoutes() {
	// Health Check
	r.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// v1 group
	v1 := r.echo.Group("/api/v1")

	// auth group
	auth := v1.Group("/auth")
	auth.POST("/register", r.handlers.User.Register)
	auth.POST("/login", r.handlers.Auth.Login)
	auth.GET("/google/login", r.handlers.Auth.GoogleLogin())
	auth.GET("/google/callback", r.handlers.Auth.GoogleCallback)
	auth.POST("/refresh", r.handlers.Auth.RefreshToken)

	protectd := v1.Group("", r.authMiddleware.Middleware)
	protectd.GET("/p", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"msg": "success"})
	})
}

func (r *Router) RegisterMinioRoutes() {
	v1 := r.echo.Group("/api/v1")
	v1.POST("/upload", r.handlers.File.Upload)
	v1.POST("/image", r.handlers.File.GetImage)
}

func (r *Router) RegisterBorrowRouter() {
	v1 := r.echo.Group("/api/v1")
	v1.POST("/borrow/return", r.handlers.Borrow.Return)
	v1.POST("/borrow", r.handlers.Borrow.Borrow)
}

func (r *Router) RegisterItemRouter() {
	v1 := r.echo.Group("/api/v1")
	v1.GET("/item/:itemID", r.handlers.Item.GetBorrowItem)
	v1.POST("/item", r.handlers.Item.CreateItem)
}
