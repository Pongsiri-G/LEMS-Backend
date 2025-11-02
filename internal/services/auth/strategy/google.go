package strategy

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	userrepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	"golang.org/x/oauth2"
)

type GoogleStrategy struct {
	client *oauth2.Config
	users  userrepo.Repository
}

func NewGoogleStrategy(client *oauth2.Config, users userrepo.Repository) *GoogleStrategy {
	return &GoogleStrategy{
		client: client,
		users:  users,
	}
}

type googleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func (s *GoogleStrategy) Authenticate(ctx context.Context, req *AuthenticateRequest) (*models.User, error) {
	if req.ProviderToken == "" {
		return nil, errors.New("missing google provider token")
	}

	// Exchange token
	token, err := s.client.Exchange(ctx, req.ProviderToken)
	if err != nil {
		return nil, err
	}

	// Get user info
	client := s.client.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var gUser googleUser
	if err := json.Unmarshal(body, &gUser); err != nil {
		return nil, err
	}

	user, err := s.users.FindByEmail(ctx, gUser.Email)

	if err != nil {
		// If not found, create new user
		newUser := &models.User{
			UserEmail:      gUser.Email,
			UserFullName:   gUser.Name,
			UserProfileURL: gUser.Picture,
			AuthProvider:   enums.Google,
			UserStatus:     enums.Pending,
		}
		if err := s.users.Create(ctx, newUser); err != nil {
			return nil, err
		}
		return nil, exceptions.ErrRegistrationSuccess
	}

	// Check user status
	switch user.UserStatus {
	case enums.Pending:
		return nil, exceptions.ErrUserPending
	case enums.Deactivated:
		return nil, exceptions.ErrUserDeactivated
	case enums.Rejected:
		return nil, exceptions.ErrUserRejected
	case enums.Active:
		return user, nil
	default:
		return nil, exceptions.ErrInactiveUser
	}
}
