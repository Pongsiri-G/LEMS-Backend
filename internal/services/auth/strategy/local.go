package strategy

import (
	"context"
	"errors"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	userrepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	"golang.org/x/crypto/bcrypt"
)

type LocalStrategy struct {
	users userrepo.Repository
}

func NewLocalStrategy(users userrepo.Repository) *LocalStrategy {
	return &LocalStrategy{users: users}
}

var ErrInvalidCredentials = errors.New("invalid email or password")

func (s *LocalStrategy) Authenticate(ctx context.Context, req *AuthenticateRequest) (*models.User, error) {
	u, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(u.UserPassword), []byte(req.Password)) != nil {
		return nil, ErrInvalidCredentials
	}
	return u, nil
}
