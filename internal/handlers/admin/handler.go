package admin

import (
	"net/http"
	"strings"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/admin"
	"github.com/labstack/echo/v4"
)

type AdminHandler interface {
	GetUsers(c echo.Context) (err error)
	GetAllUsers(c echo.Context) error

	Accept(c echo.Context) error
	Reject(c echo.Context) error
	Activate(c echo.Context) error
	Deactivate(c echo.Context) error
	Delete(c echo.Context) error
	GrantAdmin(c echo.Context) error
	RevokeAdmin(c echo.Context) error
}

type handler struct {
	adminSvc admin.AdminService
}

func NewAdminHandler(adminSvc admin.AdminService) AdminHandler {
	return &handler{adminSvc: adminSvc}

}

/* Helper functions for debugging with frontend */
func parseStatus(s string) (*enums.UserStatus, error) {
	if s == "" {
		return nil, nil
	}
	val := enums.UserStatus(strings.ToUpper(strings.TrimSpace(s)))
	switch val {
	case enums.Pending, enums.Active, enums.Rejected, enums.Deactivated:
		return &val, nil
	default:
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid status")
	}
}

func parseRole(s string) (*enums.UserRole, error) {
	if s == "" {
		return nil, nil
	}
	val := enums.UserRole(strings.ToUpper(strings.TrimSpace(s)))
	switch val {
	case enums.User, enums.Admin:
		return &val, nil
	default:
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid role")
	}
}

func (h handler) GetUsers(c echo.Context) (err error) {
	statusStr := c.QueryParam("status")
	roleStr := c.QueryParam("role")
	q := c.QueryParam("q")
	sort := c.QueryParam("sort")

	status, err := parseStatus(statusStr)
	if err != nil {
		return err
	}
	role, err := parseRole(roleStr)
	if err != nil {
		return err
	}

	var search *string
	if s := strings.TrimSpace(q); s != "" {
		search = &s
	}

	f := user.UserFilter{
		Status: status,
		Role:   role,
		Search: search,
		SortBy: sort,
	}

	users, err := h.adminSvc.GetUsers(c.Request().Context(), f)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, users)
}

func (h handler) GetAllUsers(c echo.Context) error {
	users, err := h.adminSvc.GetAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func (h handler) Accept(c echo.Context) error {
	userId := c.Param("user_id")
	err := h.adminSvc.Accept(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}

func (h handler) Reject(c echo.Context) error {
	userId := c.Param("user_id")
	err := h.adminSvc.Reject(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}

func (h handler) Activate(c echo.Context) error {
	userId := c.Param("user_id")
	err := h.adminSvc.Activate(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}

func (h handler) Deactivate(c echo.Context) error {
	userId := c.Param("user_id")
	err := h.adminSvc.Deactivate(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}

func (h handler) Delete(c echo.Context) error {
	userId := c.Param("user_id")
	err := h.adminSvc.Delete(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}

func (h handler) GrantAdmin(c echo.Context) error {
	userId := c.Param("user_id")
	err := h.adminSvc.GrantAdmin(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}

func (h handler) RevokeAdmin(c echo.Context) error {
	userId := c.Param("user_id")
	err := h.adminSvc.RevokeAdmin(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "internal server error"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}
