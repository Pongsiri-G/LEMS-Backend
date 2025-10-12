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
	logs  log.Repository
}

func newAdminService(users user.Repository, logs log.Repository) AdminService {
	return &adminService{users: users, logs: logs}
}

func (a adminService) List(ctx context.Context) ([]models.User, error) {
	return a.users.List(ctx)
}

func (a adminService) Accept(ctx context.Context, userID string) error {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	if err = a.users.UpdateStatus(ctx, userIDUUID, enums.Active); err != nil {
		return err
	}
	return nil // will save in log later
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
	if err := a.users.UpdateStatus(ctx, userIDUUID, enums.Rejected); err != nil {
		return err
	}
	return nil // will save in log later

}

func (a adminService) Deactivate(ctx context.Context, userID string) error {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	if err := a.users.UpdateStatus(ctx, userIDUUID, enums.Deactivated); err != nil {
		return err
	}
	return nil
}

func (a adminService) Delete(ctx context.Context, userID string) error {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	if err := a.users.SoftDelete(ctx, userIDUUID); err != nil {
		return err
	}
	return nil // will save in log later
}

func (a adminService) GrantAdmin(ctx context.Context, userID string) error {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	if err := a.users.UpdateRole(ctx, userIDUUID, enums.Admin); err != nil {
		return err
	}
	return nil // will save in log later
}

func (a adminService) RevokeAdmin(ctx context.Context, userID string) error {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	if err := a.users.UpdateRole(ctx, userIDUUID, enums.User); err != nil {
		return err
	}
	return nil
}
