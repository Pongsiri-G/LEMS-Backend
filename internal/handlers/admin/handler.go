package admin

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/admin"
	"github.com/labstack/echo/v4"
)

type AdminHandler interface {
	GetAllUsers(c echo.Context) error
	Accept(c echo.Context) error
	Reject(c echo.Context) error
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
