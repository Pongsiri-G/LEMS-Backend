package middlewares

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/labstack/echo/v4"
)

type rbacMiddleware struct {
	configs *configs.Config
}

type RbacMiddleware interface {
	Middleware(next echo.HandlerFunc, roles ...string) echo.HandlerFunc
}

func NewRbacMiddleware(configs *configs.Config) RbacMiddleware {
	return &rbacMiddleware{
		configs: configs,
	}
}

func (a *rbacMiddleware) Middleware(next echo.HandlerFunc, roles ...string) echo.HandlerFunc {
	return func(c echo.Context) error {
		userRole, ok := c.Get(string(enums.UserRoleContextKey)).(string)
		if !ok {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "User role not found"})
		}

		hasPermission := false

		for _, role := range roles {
			if userRole == role {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "Insufficient permission"})
		}

		return next(c)
	}
}
