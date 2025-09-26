package strategy

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	userrepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/jwt"
)

type LocalStrategy struct {
	users userrepo.Repository
	jwt   *jwt.JWTService
	cfg   *configs.Config
}

func NewLocalStrategy(users userrepo.Repository, jwt *jwt.JWTService, cfg *configs.Config) *LocalStrategy {
	return &LocalStrategy{users: users, jwt: jwt, cfg: cfg}
}

func (s *LocalStrategy) Authenticate(ctx context.Context, req *AuthenticateRequest) (*models.User, error) {
	return nil, nil
}
