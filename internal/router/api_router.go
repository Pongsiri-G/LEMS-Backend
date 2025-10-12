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

func (r *Router) RegisterAdminRoutes() {
	v1 := r.echo.Group("/api/v1")
	protected := v1.Group("", r.authMiddleware.Middleware)

	admin := protected.Group("/admin") // Maybe apply adminMiddleware later
	admin.GET("/users", r.handlers.Admin.List)
	admin.POST("/users/:id/accept", r.handlers.Admin.Accept)
	admin.POST("/users/:id/reject", r.handlers.Admin.Reject)
	admin.POST("/users/:id/deactivate", r.handlers.Admin.Deactivate)
	admin.DELETE("/users/:id", r.handlers.Admin.Delete)
	admin.POST("/users/:id/grant-admin", r.handlers.Admin.GrantAdmin)
	admin.POST("/users/:id/revoke-admin", r.handlers.Admin.RevokeAdmin)
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
	v1.GET("/item/:item-id", r.handlers.Item.GetBorrowItem)
	v1.GET("/item/list", r.handlers.Item.GetAll)
	v1.GET("/item/list/:user_id", r.handlers.Item.GetMyBorrow)
	v1.GET("/item/child/:item-id", r.handlers.Item.GetChildItemByParentID)
	v1.GET("/item/list/filter/:strategy", r.handlers.Item.GetFiltered)
	v1.POST("/item", r.handlers.Item.CreateItem)
}

func (r *Router) RegisterTagRouter() {
	v1 := r.echo.Group("/api/v1")
	v1.GET("/tag/:itemID", r.handlers.Tag.GetNameTagByItemID)
}
