package user

import (
	"context"
	"errors"
	"strings"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, r *requests.RegisterRequest) (*models.User, error)
}

type userService struct {
	userRepo user.Repository
	cfg      *configs.Config
}

func NewUserService(userRepo user.Repository, cfg *configs.Config) UserService {
	return &userService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *userService) Register(ctx context.Context, r *requests.RegisterRequest) (*models.User, error) {
	email := strings.ToLower(strings.TrimSpace(r.Email))
	phone := strings.TrimSpace(r.Phone)
	if _, err := s.userRepo.FindByEmail(ctx, r.Email); err == nil {
		return nil, exceptions.ErrEmailAlreadyExists
	} else if !errors.Is(err, user.ErrNotFound) {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		UserFullName: r.FullName,
		UserEmail:    email,
		UserPhone:    phone,
		UserPassword: string(hashed),
		UserRole:     enums.UserRole(enums.User),
		UserStatus:   enums.UserStatus(enums.Pending),
		AuthProvider: enums.AuthProvider(enums.Local),
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}
