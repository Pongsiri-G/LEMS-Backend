package router

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/ws"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/middlewares"
	"github.com/labstack/echo/v4"
)

type Router struct {
	echo           *echo.Echo
	handlers       *handlers.Handlers
	authMiddleware middlewares.AuthMiddleware
	rbacMiddleware middlewares.RbacMiddleware
	hub            *ws.Hub
}

func NewRouter(echo *echo.Echo, handlers *handlers.Handlers, authMiddleware middlewares.AuthMiddleware, rbacMiddleware middlewares.RbacMiddleware, hub *ws.Hub) *Router {
	return &Router{
		echo:           echo,
		handlers:       handlers,
		authMiddleware: authMiddleware,
		rbacMiddleware: rbacMiddleware,
		hub:            hub,
	}
}

func (r *Router) RegisterAPIRoutes() {
	// Health Check
	r.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	r.echo.GET("/ws", r.handlers.WebSocket.Run)

	// v1 group
	v1 := r.echo.Group("/api/v1")

	// auth group
	auth := v1.Group("/auth")
	auth.POST("/login", r.handlers.Auth.Login)
	auth.GET("/google/login", r.handlers.Auth.GoogleLogin())
	auth.GET("/google/callback", r.handlers.Auth.GoogleCallback)
	auth.POST("/refresh", r.handlers.Auth.RefreshToken)

	protected := v1.Group("/", r.authMiddleware.Middleware)

	// user group
	user := v1.Group("/user")
	user.POST("/register", r.handlers.User.Register)
	user.GET("/me", r.handlers.User.Me, r.authMiddleware.Middleware)

	r.registerBorrowQueueRouter(protected)

	v1.POST("/ws/noti", r.handlers.WebSocket.SendNoti, r.authMiddleware.Middleware)
}

func (r *Router) RegisterAdminRoutes() {
	v1 := r.echo.Group("/api/v1")
	protected := v1.Group("", r.authMiddleware.Middleware) // for testing

	// wrap rbac middleware to match echo.MiddlewareFunc signature
	admin := protected.Group("/admin", func(next echo.HandlerFunc) echo.HandlerFunc {
		return r.rbacMiddleware.Middleware(next, string(enums.Admin))
	})

	admin.GET("/users", r.handlers.Admin.GetAllUsers)
	admin.POST("/user/:user_id/accept", r.handlers.Admin.Accept)
	admin.POST("/user/:user_id/reject", r.handlers.Admin.Reject)
	admin.POST("/user/:user_id/activate", r.handlers.Admin.Activate)
	admin.POST("/user/:user_id/deactivate", r.handlers.Admin.Deactivate)
	admin.DELETE("/user/:user_id", r.handlers.Admin.Delete)
	admin.POST("/user/:user_id/grant-admin", r.handlers.Admin.GrantAdmin)
	admin.POST("/user/:user_id/revoke-admin", r.handlers.Admin.RevokeAdmin)
	admin.POST("/request/change-status", r.handlers.Request.ChangeRequestStatus)
	admin.GET("/logs", r.handlers.Log.GetAllLogs)
	protectd := v1.Group("", r.authMiddleware.Middleware)
	protectd.GET("/p", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"msg": "success"})
	})

	// user group
	user := protectd.Group("/user")
	user.POST("/register", r.handlers.User.Register)
	user.GET("/me", r.handlers.User.Me)
}

func (r *Router) RegisterMinioRoutes() {
	v1 := r.echo.Group("/api/v1")
	protected := v1.Group("", r.authMiddleware.Middleware)
	protected.POST("/upload", r.handlers.File.Upload)
	protected.POST("/image", r.handlers.File.GetImage)
}

func (r *Router) RegisterBorrowRouter() {
	v1 := r.echo.Group("/api/v1")
	protected := v1.Group("", r.authMiddleware.Middleware)
	protected.POST("/borrow/return", r.handlers.Borrow.Return)
	protected.POST("/borrow", r.handlers.Borrow.Borrow)
	protected.GET("/borrow/user", r.handlers.Borrow.GetMyBorrowLog)
	protected.GET("/borrow-id/:item-id", r.handlers.Borrow.GetBorrowID)
	protected.GET("/borrows", r.handlers.Borrow.GetBorrowLog)
}

func (r *Router) RegisterItemRouter() {
	v1 := r.echo.Group("/api/v1")
	protected := v1.Group("", r.authMiddleware.Middleware)
	protected.GET("/item/:item-id", r.handlers.Item.GetBorrowItem)
	protected.GET("/item/list", r.handlers.Item.GetAll)
	protected.GET("/item/list/my-borrow", r.handlers.Item.GetMyBorrow)
	protected.GET("/item/child/:item-id", r.handlers.Item.GetChildItemByParentID)
	protected.GET("/item/list/search", r.handlers.Item.SearchItems)
	protected.POST("/item", r.handlers.Item.CreateItem)
	protected.DELETE("/item/:item-id", r.handlers.Item.DeleteItem)

	protected.PUT("/item", r.handlers.Item.EditItem)
	protected.POST("/item/assign-item-set/:parent_id/:child_id", r.handlers.Item.AssignItemSet)
	protected.POST("/item/remove-item-set/:parent_id/:child_id", r.handlers.Item.RemoveItemSet)
}

func (r *Router) RegisterTagRouter() {
	v1 := r.echo.Group("/api/v1")
	protected := v1.Group("", r.authMiddleware.Middleware)
	protected.POST("/tag", r.handlers.Tag.CreateTag)
	protected.GET("/tags", r.handlers.Tag.GetTags)
	protected.GET("/tag/:itemID", r.handlers.Tag.GetNameTagByItemID)
	protected.DELETE("/tag/unassign/:item_id/:tag_id", r.handlers.Tag.UnAssignTagFromItem)
	protected.POST("/tag/assign/:item_id/:tag_id", r.handlers.Tag.AssignTagToItem)
}

func (r *Router) RegisterRequestRouter() {
	v1 := r.echo.Group("/api/v1")
	protected := v1.Group("", r.authMiddleware.Middleware)
	protected.GET("/requests", r.handlers.Request.GetRequests)
	protected.GET("/requests/user", r.handlers.Request.GetMyRequests)
	protected.POST("/request", r.handlers.Request.CreateRequest)
	protected.PUT("/request", r.handlers.Request.EditRequest)
	protected.POST("/request/cancel", r.handlers.Request.CancelRequest)
	protected.POST("/requests/export", r.handlers.Request.ExportRequests)
}

func (r *Router) registerBorrowQueueRouter(route *echo.Group) {
	bq := route.Group("bq")
	bq.POST("/enqueue", r.handlers.BorrowQueue.Enqueue)
	bq.GET("/front", r.handlers.BorrowQueue.GetFrontQueue)
}
