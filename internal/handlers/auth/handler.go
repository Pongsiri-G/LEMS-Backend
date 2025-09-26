package auth

import (
	"net/http"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// ---- Dummy endpoints (Phase 1) ----
func (h *AuthHandler) Register(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]string{"message": "dummy register"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "dummy login"})
}

func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "dummy google login"})
}

func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "dummy google callback"})
}
