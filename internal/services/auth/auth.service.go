package services

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	userrepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth/strategy"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/jwt"
)

type AuthService struct {
	strategies map[string]strategy.AuthStrategy
	users      userrepo.Repository
	jwt        *jwt.JWTService
}

func NewAuthService(strategies map[string]strategy.AuthStrategy, users userrepo.Repository, jwt *jwt.JWTService) *AuthService {
	return &AuthService{strategies: strategies, users: users, jwt: jwt}
}

func (s *AuthService) Login(ctx context.Context, key string, req *strategy.AuthenticateRequest) (*models.User, error) {
	return nil, nil
}

func (s *AuthService) Register(ctx context.Context, u *models.User) (*models.User, error) {
	return nil, nil
}
