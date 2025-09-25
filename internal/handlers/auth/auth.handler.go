package handlers

import (
	"context"

	services "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth"
)

type AuthHandler interface {
	LocalSignIn(ctx context.Context, username string, password string) error
}

type authHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

// LocalSignIn implements AuthHandler.
func (a *authHandler) LocalSignIn(ctx context.Context, username string, password string) error {
	panic("unimplemented")
}
