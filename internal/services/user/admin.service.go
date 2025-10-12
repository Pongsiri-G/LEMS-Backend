package user

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/log"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	"github.com/google/uuid"
)

type AdminService interface {
	List(ctx context.Context) ([]models.User, error)

	Accept(ctx context.Context, userID string) error
	Reject(ctx context.Context, userID string) error
	Deactivate(ctx context.Context, userID string) error
	Delete(ctx context.Context, userID string) error
	GrantAdmin(ctx context.Context, userID string) error
	RevokeAdmin(ctx context.Context, userID string) error
}

type adminService struct {
	users user.Repository
}

func newAdminService(users user.Repository, logs log.Repository) AdminService {
	return &adminService{users: users}
}

func (a adminService) List(ctx context.Context) ([]models.User, error) {
	return a.users.List(ctx)
}

func (a adminService) Accept(ctx context.Context, userID string) error {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return a.users.UpdateStatus(ctx, userIDUUID, enums.Active)
}

func (a adminService) Reject(ctx context.Context, userID string) error {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	// check if user is pending
	_, err = a.users.FindById(ctx, userIDUUID)
	if err != nil {
		return user.ErrRejectOnlyPending // Rejected only pending user
	}
	return a.users.UpdateStatus(ctx, userIDUUID, enums.Rejected)
}

func (a adminService) Deactivate(ctx context.Context, userID string) error {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return a.users.UpdateStatus(ctx, userIDUUID, enums.Deactivated)
}

func (a adminService) Delete(ctx context.Context, userID string) error {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return a.users.SoftDelete(ctx, userIDUUID)
}

func (a adminService) GrantAdmin(ctx context.Context, userID string) error {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return a.users.UpdateRole(ctx, userIDUUID, enums.Admin)
}

func (a adminService) RevokeAdmin(ctx context.Context, userID string) error {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return a.users.UpdateRole(ctx, userIDUUID, enums.User)
}
