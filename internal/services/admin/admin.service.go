package admin

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	"github.com/google/uuid"
)

type AdminService interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)

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

func NewAdminService(users user.Repository) AdminService {
	return &adminService{users: users}
}

func (a adminService) checkStatus(u *models.User) error {
	if u.UserStatus != enums.Pending {
		return user.ErrUserIsNotPending
	}
	return nil
}

func (a adminService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return a.users.GetAllUsers(ctx)
}

func (a adminService) Accept(ctx context.Context, userID string) error {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	u, err := a.users.FindById(ctx, userIDUUID)
	if err != nil {
		return err
	}

	err = a.checkStatus(u)
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

	u, err := a.users.FindById(ctx, userIDUUID)
	if err != nil {
		return err
	}

	err = a.checkStatus(u)
	if err != nil {
		return err
	}

	return a.users.UpdateStatus(ctx, userIDUUID, enums.Rejected)
}

func (a adminService) Deactivate(ctx context.Context, userID string) error {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	u, err := a.users.FindById(ctx, userIDUUID)
	if err != nil {
		return err
	}

	if u.UserStatus == enums.Pending {
		return user.ErrDeactivatePending
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

	u, err := a.users.FindById(ctx, userIDUUID)
	if err != nil {
		return err
	}

	if u.UserRole == enums.Admin {
		return user.ErrAlreadyAdmin
	}

	return a.users.UpdateRole(ctx, userIDUUID, enums.Admin)
}

func (a adminService) RevokeAdmin(ctx context.Context, userID string) error {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	u, err := a.users.FindById(ctx, userIDUUID)
	if err != nil {
		return err
	}

	if u.UserRole == enums.Admin {
		return user.ErrRevokeUser
	}
	return a.users.UpdateRole(ctx, userIDUUID, enums.User)
}
