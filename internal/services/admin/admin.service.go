package admin

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	logrepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/log"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type AdminService interface {
	// --- Query --
	GetUser(ctx context.Context, userID string) (*models.User, error)
	GetUsers(ctx context.Context, filter user.UserFilter) ([]models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)

	// --- Command ---
	Accept(ctx context.Context, adminID, userID string) error
	Reject(ctx context.Context, adminID, userID string) error
	Activate(ctx context.Context, adminID, userID string) error
	Deactivate(ctx context.Context, adminID, userID string) error
	Delete(ctx context.Context, adminID, userID string) error
	GrantAdmin(ctx context.Context, adminID, userID string) error
	RevokeAdmin(ctx context.Context, adminID, userID string) error
}

type adminService struct {
	users   user.Repository
	logRepo logrepo.Repository
}

func NewAdminService(users user.Repository, logRepo logrepo.Repository) AdminService {
	return &adminService{
		users:   users,
		logRepo: logRepo,
	}
}

func (a adminService) checkPending(u *models.User) error {
	if u.UserStatus != enums.Pending {
		return user.ErrUserIsNotPending
	}
	return nil
}

func (a adminService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	u, err := a.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (a adminService) GetUsers(ctx context.Context, filter user.UserFilter) ([]models.User, error) {
	return a.users.List(ctx, filter)
}

func (a adminService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return a.users.GetAllUsers(ctx)
}

func (a adminService) Accept(ctx context.Context, adminID, userID string) error {
	u, err := a.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := a.checkPending(u); err != nil {
		return err
	}
	if err := a.users.UpdateStatus(ctx, u.UserID, enums.Active); err != nil {
		return err
	}

	// Log admin action
	adminUUID, _ := uuid.Parse(adminID)
	if err := a.logRepo.CreateAdminActionLog(ctx, adminUUID, enums.LogTypeAccept, u.UserID); err != nil {
		log.Error().Err(err).Msg("Failed to create admin action log")
	}

	return nil
}

func (a adminService) Reject(ctx context.Context, adminID, userID string) error {
	u, err := a.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := a.checkPending(u); err != nil {
		return err
	}
	if err := a.users.UpdateStatus(ctx, u.UserID, enums.Rejected); err != nil {
		return err
	}

	// Log admin action
	adminUUID, _ := uuid.Parse(adminID)
	if err := a.logRepo.CreateAdminActionLog(ctx, adminUUID, enums.LogTypeReject, u.UserID); err != nil {
		log.Error().Err(err).Msg("Failed to create admin action log")
	}

	return nil
}

func (a adminService) Activate(ctx context.Context, adminID, userID string) error {
	u, err := a.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if u.UserStatus != enums.Deactivated {
		return user.ErrAlreadyActiveOrStillPending
	}
	if err := a.users.UpdateStatus(ctx, u.UserID, enums.Active); err != nil {
		return err
	}

	// Log admin action
	adminUUID, _ := uuid.Parse(adminID)
	if err := a.logRepo.CreateAdminActionLog(ctx, adminUUID, enums.LogTypeActivate, u.UserID); err != nil {
		log.Error().Err(err).Msg("Failed to create admin action log")
	}

	return nil
}

func (a adminService) Deactivate(ctx context.Context, adminID, userID string) error {
	u, err := a.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if u.UserStatus == enums.Pending {
		return user.ErrDeactivatePending
	}
	if err := a.users.UpdateStatus(ctx, u.UserID, enums.Deactivated); err != nil {
		return err
	}

	// Log admin action
	adminUUID, _ := uuid.Parse(adminID)
	if err := a.logRepo.CreateAdminActionLog(ctx, adminUUID, enums.LogTypeDeactivate, u.UserID); err != nil {
		log.Error().Err(err).Msg("Failed to create admin action log")
	}

	return nil
}

func (a adminService) Delete(ctx context.Context, adminID, userID string) error {
	u, err := a.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := a.users.SoftDelete(ctx, u.UserID); err != nil {
		return err
	}

	// Log admin action
	adminUUID, _ := uuid.Parse(adminID)
	if err := a.logRepo.CreateAdminActionLog(ctx, adminUUID, enums.LogTypeDelete, u.UserID); err != nil {
		log.Error().Err(err).Msg("Failed to create admin action log")
	}

	return nil
}

func (a adminService) GrantAdmin(ctx context.Context, adminID, userID string) error {
	u, err := a.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if u.UserRole == enums.Admin {
		return user.ErrAlreadyAdmin
	}
	if err := a.users.UpdateRole(ctx, u.UserID, enums.Admin); err != nil {
		return err
	}

	// Log admin action
	adminUUID, _ := uuid.Parse(adminID)
	if err := a.logRepo.CreateAdminActionLog(ctx, adminUUID, enums.LogTypeGrantAdmin, u.UserID); err != nil {
		log.Error().Err(err).Msg("Failed to create admin action log")
	}

	return nil
}

func (a adminService) RevokeAdmin(ctx context.Context, adminID, userID string) error {
	u, err := a.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if u.UserRole != enums.Admin {
		return user.ErrRevokeUser
	}
	if err := a.users.UpdateRole(ctx, u.UserID, enums.User); err != nil {
		return err
	}

	// Log admin action
	adminUUID, _ := uuid.Parse(adminID)
	if err := a.logRepo.CreateAdminActionLog(ctx, adminUUID, enums.LogTypeRevokeAdmin, u.UserID); err != nil {
		log.Error().Err(err).Msg("Failed to create admin action log")
	}

	return nil
}
