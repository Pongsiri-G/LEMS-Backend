package user

import (
	"context"
	"errors"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	logrepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/log"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, r *requests.RegisterRequest) (*responses.UserResponse, error)
	FindByID(ctx context.Context, userID string) (*responses.UserResponse, error)
}

type userService struct {
	userRepo user.Repository
	logRepo  logrepo.Repository
	cfg      *configs.Config
}

func NewUserService(userRepo user.Repository, logRepo logrepo.Repository, cfg *configs.Config) UserService {
	return &userService{
		userRepo: userRepo,
		logRepo:  logRepo,
		cfg:      cfg,
	}
}

func (s *userService) Register(ctx context.Context, r *requests.RegisterRequest) (*responses.UserResponse, error) {
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
		UserEmail:    r.Email,
		UserPhone:    r.Phone,
		UserPassword: string(hashed),
		UserRole:     enums.UserRole(enums.User),
		UserStatus:   enums.UserStatus(enums.Pending),
		AuthProvider: enums.AuthProvider(enums.Local),
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	// Create register log
	if err := s.logRepo.CreateRegisterLog(ctx, u.UserID); err != nil {
		log.Error().Err(err).Msg("Failed to create register log")
	}

	return s.toResponse(u), nil
}

// MyInfo implements UserService.
func (s *userService) FindByID(ctx context.Context, userID string) (*responses.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(user), nil
}

func (s *userService) toResponse(u *models.User) *responses.UserResponse {
	return &responses.UserResponse{
		UserID:         u.UserID.String(),
		UserFullName:   u.UserFullName,
		UserEmail:      u.UserEmail,
		UserPhone:      u.UserPhone,
		UserRole:       u.UserRole,
		UserStatus:     u.UserStatus,
		UserProfileURL: u.UserProfileURL,
		AuthProvider:   u.AuthProvider,
		LastLoggedIn:   u.LastLoggedIn,
	}
}
