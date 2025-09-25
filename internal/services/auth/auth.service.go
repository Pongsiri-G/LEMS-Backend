package services

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	repositories "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
)

type AuthService interface {
	LocalSignIn(ctx context.Context, username string, password string) error
	GoogleSignIn(ctx context.Context, gmail string) error
}

type authService struct {
	repo   repositories.UserRepoistory
	config *configs.Config
}

func NewAuthService(
	config *configs.Config,
	repo repositories.UserRepoistory,
) AuthService {
	return &authService{
		config: config,
		repo:   repo,
	}
}

// GoogleSignIn implements AuthService.
func (a *authService) GoogleSignIn(ctx context.Context, gmail string) error {
	panic("unimplemented")
}

// LocalSignIn implements AuthService.
func (a *authService) LocalSignIn(ctx context.Context, username string, password string) error {
	panic("unimplemented")
}
