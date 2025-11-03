package middlewares

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/user"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type rbacMiddleware struct {
	configs *configs.Config
	userService user.UserService
}

type RbacMiddleware interface {
	Middleware(next echo.HandlerFunc, roles ...enums.UserRole) echo.HandlerFunc
}

func NewRbacMiddleware(configs *configs.Config, userService user.UserService) RbacMiddleware {
	return &rbacMiddleware{
		configs: configs,
		userService: userService,
	}
}

func (a *rbacMiddleware) Middleware(next echo.HandlerFunc, roles ...enums.UserRole) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, ok := c.Get(string(enums.UserIDContextKey)).(string)
		if !ok {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "User id not found"})
		}

		user, err := a.userService.FindByID(c.Request().Context(), userID)
		if err != nil {
			log.Err(err).Msg("failed to find user")
			return c.JSON(http.StatusForbidden, map[string]string{"message": err.Error()})
		}

		hasPermission := false

		for _, role := range roles {
			if user.UserRole == role {
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
