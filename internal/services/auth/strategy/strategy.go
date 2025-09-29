package strategy

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
)

type AuthenticateRequest struct {
	// Local
	Email    string
	Password string
	// Google
	ProviderToken string
}

type AuthStrategy interface {
	Authenticate(ctx context.Context, req *AuthenticateRequest) (*models.User, error)
}
