package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth/strategy"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

type AuthHandler interface {
	Login(c echo.Context) error
	GoogleLogin() echo.HandlerFunc
	GoogleCallback(c echo.Context) error
	RefreshToken(c echo.Context) error
}

type authHandler struct {
	svc   auth.AuthService
	oauth *oauth2.Config
	cfg   *configs.Config
}

func NewAuthHandler(svc auth.AuthService, oauth *oauth2.Config, cfg *configs.Config) AuthHandler {
	return &authHandler{
		svc:   svc,
		oauth: oauth,
		cfg:   cfg,
	}
}

func (h *authHandler) Login(c echo.Context) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&body); err != nil || body.Email == "" || body.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	res, err := h.svc.Login(
		c.Request().Context(),
		"local",
		&strategy.AuthenticateRequest{Email: body.Email, Password: body.Password})
	if err != nil {
		log.Err(err)

		switch err {
		case exceptions.ErrUserPending, exceptions.ErrUserDeactivated, exceptions.ErrUserRejected, exceptions.ErrInactiveUser:
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, map[string]any{"access_token": res.AccessToken, "refresh_token": res.RefreshToken})
}

// GET /api/v1/auth/google/login
func (h *authHandler) GoogleLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		url := h.oauth.AuthCodeURL("random-state-string", oauth2.AccessTypeOffline)
		return c.JSON(http.StatusOK, url)
	}
}

// GET /api/v1/auth/google/callback?code=...
func (h *authHandler) GoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")

	res, err := h.svc.Login(
		c.Request().Context(),
		"google",
		&strategy.AuthenticateRequest{
			ProviderToken: code,
		},
	)

	if err != nil {
		log.Err(err).Send()

		switch err {
		case exceptions.ErrRegistrationSuccess, exceptions.ErrUserPending, exceptions.ErrUserDeactivated, exceptions.ErrUserRejected, exceptions.ErrInactiveUser:
			return c.Redirect(http.StatusTemporaryRedirect, h.redirectWithMsg(err.Error()))
		default:
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
	}

	return c.Redirect(http.StatusTemporaryRedirect, h.redirectWithTokens(res))
}

func (h *authHandler) RefreshToken(c echo.Context) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	res, err := h.svc.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(200, map[string]interface{}{
		"access_token": res.AccessToken,
		"user":         res.User,
	})
}

func (h *authHandler) redirectWithTokens(response *responses.AuthResponse) string {
	frontendURL := fmt.Sprintf("%s/oauth/callback", h.cfg.AllowOrigins[0])
	fmt.Print(frontendURL)

	resURL := fmt.Sprintf(
		"%s?accessToken=%s&refreshToken=%s",
		frontendURL,
		url.QueryEscape(response.AccessToken),
		url.QueryEscape(response.RefreshToken),
	)

	return resURL
}

func (h *authHandler) redirectWithMsg(massage string) string {
	frontendURL := fmt.Sprintf("%s/oauth/callback", h.cfg.AllowOrigins[0])
	fmt.Print(frontendURL)

	resURL := fmt.Sprintf(
		"%s?msg=%s",
		frontendURL,
		url.QueryEscape(massage),
	)

	return resURL
}
