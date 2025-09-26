package strategy

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/google"
	userrepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
)

type GoogleStrategy struct {
	client *google.GoogleOAuthClient
	users  userrepo.Repository
	cfg    *configs.Config
}

func NewGoogleStrategy(client *google.GoogleOAuthClient, users userrepo.Repository, cfg *configs.Config) *GoogleStrategy {
	return &GoogleStrategy{client: client, users: users, cfg: cfg}
}

func (s *GoogleStrategy) Authenticate(ctx context.Context, req *AuthenticateRequest) (*models.User, error) {
	return nil, nil
}
