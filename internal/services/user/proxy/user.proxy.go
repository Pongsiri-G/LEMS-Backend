package proxy

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/user"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils/contextutil"
	"github.com/labstack/echo/v4"
)

type UserServiceProxy struct {
	service user.UserService
}

func NewUserServiceProxy(real *user.UserService) *UserServiceProxy {
	return &UserServiceProxy{
		service: *real,
	}
}

// Delete implements UserService.
func (u *UserServiceProxy) Delete(c echo.Context, targetID string) error {
	authUser, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return exceptions.ErrNotAllowAccess
	}

	if isAdmin(enums.UserRole(authUser.Role)) {
		return exceptions.ErrNotAllowAccess
	}

	// return u.service.Delete(c.Request().Context(), targetID)
	return nil
}

// MyInfo implements UserService.
func (u *UserServiceProxy) FindByID(ctx context.Context, userID string) (*responses.UserResponse, error) {
	return u.service.FindByID(ctx, userID)
}

// Register implements UserService.
func (u *UserServiceProxy) Register(ctx context.Context, r *requests.RegisterRequest) (*responses.UserResponse, error) {
	return u.service.Register(ctx, r)
}

// UpdateUserStatus implements UserService.
func (u *UserServiceProxy) UpdateUserStatus(c echo.Context, targetID string, status string) (*responses.UserResponse, error) {
	authUser, err := contextutil.GetUserFromContext(c)
	if err != nil {
		return nil, exceptions.ErrNotAllowAccess
	}

	if isAdmin(enums.UserRole(authUser.Role)) {
		return nil, exceptions.ErrNotAllowAccess
	}

	// return u.service.UpdateUserStatus(c.Request().Context(), targetID, status)
	panic("unimplement")
}

func isAdmin(role enums.UserRole) bool {
	return role == enums.UserRole(enums.Admin)
}
