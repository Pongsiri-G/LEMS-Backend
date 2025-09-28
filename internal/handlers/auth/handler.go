package auth

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth/strategy"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type AuthHandler interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
	GoogleLogin() echo.HandlerFunc
	GoogleCallback(c echo.Context) error
}

type authHandler struct {
	svc   services.AuthService
	oauth *oauth2.Config
}

func NewAuthHandler(svc services.AuthService, oauth *oauth2.Config) AuthHandler {
	return &authHandler{
		svc:   svc,
		oauth: oauth,
	}
}

// POST /api/v1/auth/register (LOCAL)
func (h *authHandler) Register(c echo.Context) error {
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

func (h *authHandler) Login(c echo.Context) error {
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
func (h *authHandler) GoogleLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		url := h.oauth.AuthCodeURL("random-state-string", oauth2.AccessTypeOffline)
		return c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

// GET /api/v1/auth/google/callback?code=...
func (h *authHandler) GoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")

	user, token, err := h.svc.Login(
		c.Request().Context(),
		"google",
		&strategy.AuthenticateRequest{
			ProviderToken: code,
		},
	)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	return c.JSON(
		http.StatusOK,
		map[string]any{
			"user":  user,
			"token": token,
		},
	)

}
