package middlewares

import (
	"fmt"
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/user"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils/tokenutil"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type authMiddleware struct {
	configs *configs.Config
	userService user.UserService
}

type AuthMiddleware interface {
	Middleware(next echo.HandlerFunc) echo.HandlerFunc
}

func NewAuthMiddleware(configs *configs.Config, userService user.UserService) AuthMiddleware {
	return &authMiddleware{
		configs: configs,
		userService: userService,
	}
}

func (a *authMiddleware) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString, err := tokenutil.GetTokenFromEchoHeader(c)
		if err != nil {
			log.Err(err).Msg("error")
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
		}

		if err != nil || tokenString == "" {
			log.Err(err).Msg("error")
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
		}
		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &auth.JWTClaims{}, func(token *jwt.Token) (any, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Err(err).Send()
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(a.configs.JWT.JwtSecret), nil
		})
		if err != nil {
			log.Err(err).Msgf("69: %s", tokenString)
			return c.JSON(http.StatusForbidden, map[string]string{"message": err.Error()})
		}

		// Validate claims
		claims, ok := token.Claims.(*auth.JWTClaims)
		if !ok || !token.Valid {
			log.Err(err).Msg("76")
			return c.JSON(http.StatusForbidden, map[string]string{"message": "invalid claim"})
		}

		// Validate current user
		user, err := a.userService.FindByID(c.Request().Context(), claims.UserID)
		if err != nil {
			log.Err(err).Msg("failed to find user")
			return c.JSON(http.StatusForbidden, map[string]string{"message": err.Error()})
		}

		if user.UserStatus != enums.Active {
			return c.JSON(http.StatusForbidden, map[string]string{"message": exceptions.ErrInactiveUser.Error()})
		}

		// Set claims and user ID to context
		c.Set(string(enums.UserIDContextKey), claims.UserID)
		c.Set(string(enums.UserEmailContextKey), claims.Email)
		c.Set(string(enums.UserRoleContextKey), claims.Role)

		return next(c)
	}
}
