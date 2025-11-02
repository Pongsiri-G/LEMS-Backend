package user

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/user"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils/contextutil"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type UserHandler interface {
	Register(c echo.Context) error
	Me(c echo.Context) error
}

type userHandler struct {
	oauth     *oauth2.Config
	userService user.UserService
}

func NewUserHandler(userService user.UserService, oauth *oauth2.Config) UserHandler {
	return &userHandler{
		oauth:     oauth,
		userService: userService, 
	}
}

// POST /api/v1/auth/register (LOCAL)
func (h *userHandler) Register(c echo.Context) error {
	var body struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	if err := c.Bind(&body); err != nil || body.Email == "" || body.Password == "" || body.FullName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	// register พร้อมทั้งอัปเดทไปยังฐานข้อมูลหากลงทะเบียนสำเร็จ
	_, err := h.userService.Register(c.Request().Context(), &requests.RegisterRequest{
		FullName: body.FullName,
		Email:    body.Email,
		Phone:    body.Phone,
		Password: body.Password,
	})
	if err != nil {
		if err == exceptions.ErrEmailAlreadyExists {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "email already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": exceptions.ErrInternalServer.Error()})
	}
	// สมัครเสร็จให้ไป login ต่อ
	return c.JSON(http.StatusCreated, map[string]string{"message": "register success"})
}

func (h *userHandler) Me(c echo.Context) error {
	authUser, err := contextutil.GetUserFromContext(c)

	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}

	res, err := h.userService.FindByID(
		c.Request().Context(),
		authUser.ID,
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}
