package tokenutil

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

func SplitBearerToken(bearer string) (string, error) {
	bearer = strings.TrimSpace(bearer)
	splittedToken := strings.Split(bearer, "Bearer ")
	if len(splittedToken) != 2 {
		return "", ErrInvalidToken
	}

	token := splittedToken[1]

	return token, nil
}

func GetTokenFromEchoHeader(c echo.Context) (string, error) {
	bearer := c.Request().Header.Get("Authorization")

	token, err := SplitBearerToken(bearer)
	if err != nil {
		return "", err
	}

	return token, nil
}

func GetTokenFromEchoUrl(c echo.Context) (string, error) {
	query := c.Request().URL.Query()
	bearer := query.Get("accTk")

	token, err := SplitBearerToken(bearer)
	if err != nil {
		return "", err
	}

	return token, nil
}
