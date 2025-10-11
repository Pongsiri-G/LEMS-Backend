package middlewares

import (
	"fmt"
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils/tokenutil"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type authMiddleware struct {
	configs *configs.Config
}

type AuthMiddleware interface {
	Middleware(next echo.HandlerFunc) echo.HandlerFunc
}

func NewAuthMiddleware(configs *configs.Config) AuthMiddleware {
	return &authMiddleware{
		configs: configs,
	}
}

func (a *authMiddleware) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString, err := tokenutil.GetTokenFromEchoHeader(c)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"massage": err.Error()})
		}
		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &auth.JWTClaims{}, func(token *jwt.Token) (any, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(a.configs.JWT.JwtSecret), nil
		})
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusForbidden, map[string]string{"massage": err.Error()})
		}

		// Validate claims
		claims, ok := token.Claims.(*auth.JWTClaims)
		if !ok || !token.Valid {
			return c.JSON(http.StatusForbidden, map[string]string{"massage": err.Error()})
		}

		// Set claims and user ID to context
		c.Set(string(enums.UserIDContextKey), claims.UserID)

		return next(c)
	}
}
