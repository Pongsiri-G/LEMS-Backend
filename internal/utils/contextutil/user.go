package contextutil

import (
	"errors"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AuthUser struct {
	ID    string
	Email string
	Role  string
}

func GetUserFromContext(c echo.Context) (*AuthUser, error) {
	id, ok := c.Get(string(enums.UserIDContextKey)).(string)
	if !ok || id == "" {
		log.Info().Msgf("msg: %s", id)
		return nil, errors.New("user id not found in context")
	}

	email, ok := c.Get(string(enums.UserEmailContextKey)).(string)
	if !ok || email == "" {
		log.Info().Msgf("msg: %s", id)
		return nil, errors.New("email not found in context")
	}

	role, ok := c.Get(string(enums.UserRoleContextKey)).(string)
	if !ok || role == "" {
		log.Info().Msgf("msg: %s", id)
		return nil, errors.New("role not found in context")
	}

	return &AuthUser{
		ID:    id,
		Email: email,
		Role:  role,
	}, nil
}
