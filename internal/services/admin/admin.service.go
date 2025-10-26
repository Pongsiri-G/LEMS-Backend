package admin

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	"github.com/google/uuid"
)

type AdminService interface {
	// --- Query --
	GetUser(ctx context.Context, userID string) (*models.User, error)
	GetUsers(ctx context.Context, filter user.UserFilter) ([]models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)

	// --- Command ---
	Accept(ctx context.Context, userID string) error
	Reject(ctx context.Context, userID string) error
	Activate(ctx context.Context, userID string) error
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

func (a adminService) parseAndGet(ctx context.Context, userID string) (*models.User, uuid.UUID, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, uuid.Nil, err
	}
	u, err := a.users.FindById(ctx, uid)
	if err != nil {
		return nil, uuid.Nil, err
	}
	return u, uid, nil
}

func (a adminService) checkPending(u *models.User) error {
	if u.UserStatus != enums.Pending {
		return user.ErrUserIsNotPending
	}
	return nil
}

func (a adminService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	_, uid, err := a.parseAndGet(ctx, userID)
	if err != nil {
		return nil, err
	}
	return a.users.FindById(ctx, uid)
}

func (a adminService) GetUsers(ctx context.Context, filter user.UserFilter) ([]models.User, error) {
	return a.users.List(ctx, filter)
}

func (a adminService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return a.users.GetAllUsers(ctx)
}

func (a adminService) Accept(ctx context.Context, userID string) error {
	u, uid, err := a.parseAndGet(ctx, userID)
	if err != nil {
		return err
	}
	if err := a.checkPending(u); err != nil {
		return err
	}
	return a.users.UpdateStatus(ctx, uid, enums.Active)
}

func (a adminService) Reject(ctx context.Context, userID string) error {
	u, uid, err := a.parseAndGet(ctx, userID)
	if err != nil {
		return err
	}
	if err := a.checkPending(u); err != nil {
		return err
	}
	return a.users.UpdateStatus(ctx, uid, enums.Rejected)
}

func (a adminService) Activate(ctx context.Context, userID string) error {
	u, uid, err := a.parseAndGet(ctx, userID)
	if err != nil {
		return err
	}
	if u.UserStatus != enums.Deactivated {
		return user.ErrAlreadyActiveOrStillPending
	}
	return a.users.UpdateStatus(ctx, uid, enums.Active)
}

func (a adminService) Deactivate(ctx context.Context, userID string) error {
	u, uid, err := a.parseAndGet(ctx, userID)
	if err != nil {
		return err
	}
	if u.UserStatus == enums.Pending {
		return user.ErrDeactivatePending
	}
	return a.users.UpdateStatus(ctx, uid, enums.Deactivated)
}

func (a adminService) Delete(ctx context.Context, userID string) error {
	_, uid, err := a.parseAndGet(ctx, userID)
	if err != nil {
		return err
	}
	return a.users.SoftDelete(ctx, uid)
}

func (a adminService) GrantAdmin(ctx context.Context, userID string) error {
	u, uid, err := a.parseAndGet(ctx, userID)
	if err != nil {
		return err
	}
	if u.UserRole == enums.Admin {
		return user.ErrAlreadyAdmin
	}
	return a.users.UpdateRole(ctx, uid, enums.Admin)
}

func (a adminService) RevokeAdmin(ctx context.Context, userID string) error {
	u, uid, err := a.parseAndGet(ctx, userID)
	if err != nil {
		return err
	}

	if u.UserRole != enums.Admin {
		return user.ErrRevokeUser
	}
	return a.users.UpdateRole(ctx, uid, enums.User)
}
