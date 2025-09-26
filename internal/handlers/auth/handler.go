package auth

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth/strategy"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// POST /api/v1/auth/register (LOCAL)
func (h *AuthHandler) Register(c echo.Context) error {
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
	_, err := h.svc.Register(c.Request().Context(), &services.RegisterRequest{
		FullName: body.FullName,
		Email:    body.Email,
		Phone:    body.Phone,
		Password: body.Password,
	})
	if err != nil {
		if err == services.ErrEmailAlreadyExists {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "email already exists"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	// สมัครเสร็จให้ไป login ต่อ
	return c.JSON(http.StatusCreated, map[string]string{"message": "register success"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&body); err != nil || body.Email == "" || body.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	u, token, err := h.svc.Login(c.Request().Context(),
		"local",
		&strategy.AuthenticateRequest{Email: body.Email, Password: body.Password})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	return c.JSON(http.StatusOK, map[string]any{"user": u, "token": token})
}

// GET /api/v1/auth/google/login
func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "dummy google login"})
}

// GET /api/v1/auth/google/callback?code=...
func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "dummy google callback"})
}
