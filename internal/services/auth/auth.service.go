package services

import (
	"context"
	"errors"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	userrepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth/strategy"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, r *RegisterRequest) (*models.User, error)
	Login(ctx context.Context, key string, req *strategy.AuthenticateRequest) (*models.User, string, error)
}

type authService struct {
	strategies map[string]strategy.AuthStrategy
	users      userrepo.Repository
	jwt        *jwt.JWTService
}

func NewAuthService(strategies map[string]strategy.AuthStrategy, users userrepo.Repository, jwt *jwt.JWTService) AuthService {
	return &authService{
		strategies: strategies,
		users:      users,
		jwt:        jwt,
	}
}

func (s *authService) Login(ctx context.Context, key string, req *strategy.AuthenticateRequest) (*models.User, string, error) {
	strategy, ok := s.strategies[key]
	if !ok {
		return nil, "", errors.New("strategy not found")
	}
	u, err := strategy.Authenticate(ctx, req)
	if err != nil {
		return nil, "", err
	}
	token, err := s.jwt.Generate(u.UserID)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}

type RegisterRequest struct {
	FullName string
	Email    string
	Password string
	Phone    string
}

var ErrEmailAlreadyExists = errors.New("email already exists")

func (s *authService) Register(ctx context.Context, r *RegisterRequest) (*models.User, error) {
	if _, err := s.users.FindByEmail(ctx, r.Email); err == nil {
		return nil, ErrEmailAlreadyExists
	} else if !errors.Is(err, userrepo.ErrNotFound) {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		UserFullName: r.FullName,
		UserEmail:    r.Email,
		UserPhone:    r.Phone,
		UserPassword: string(hashed),
		UserRole:     enums.UsesrRole("USER"),
		UserStatus:   enums.UserStatus("PENDING"),
		AuthProvider: enums.AuthProvider("LOCAL"),
	}

	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}
