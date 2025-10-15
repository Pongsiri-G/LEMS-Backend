package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/auth"
	userrepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/auth/strategy"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils/timeutil"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type AuthService interface {
	Login(ctx context.Context, key string, req *strategy.AuthenticateRequest) (*responses.AuthResponse, error)
	RefreshToken(ctx context.Context, tokenStr string) (*responses.AuthResponse, error)
}

type authService struct {
	cfg        *configs.Config
	strategies map[string]strategy.AuthStrategy
	users      userrepo.Repository
}

func NewAuthService(strategies map[string]strategy.AuthStrategy, users userrepo.Repository, cfg *configs.Config) AuthService {
	return &authService{
		strategies: strategies,
		users:      users,
		cfg:        cfg,
	}
}

func (s *authService) Login(ctx context.Context, key string, req *strategy.AuthenticateRequest) (*responses.AuthResponse, error) {
	strategy, ok := s.strategies[key]
	if !ok {
		return nil, errors.New("strategy not found")
	}
	u, err := strategy.Authenticate(ctx, req)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(err.Error())
	}

	accessToken, refreshToken, err := s.generateJWTToken(u.UserID.String(), u.UserEmail)

	if err != nil {
		log.Debug().Stack()
		return nil, err
	}

	// Update last login
	s.users.UpdateLastLogin(ctx, u.UserID)

	return &responses.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, tokenStr string) (*responses.AuthResponse, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &auth.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWT.RefreshSecret), nil
	})

	if err != nil {
		log.Err(err).Msg("")
		return nil, err
	}

	claims, ok := token.Claims.(*auth.JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	access, refresh, err := s.generateJWTToken(claims.UserID, claims.Email)
	if err != nil {
		return nil, err
	}

	return &responses.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *authService) generateJWTToken(userID string, userEmail string) (string, string, error) {
	now := timeutil.BangkokNow()

	accessTokenExpiration, err := time.ParseDuration(s.cfg.JWT.JwtExpirationMinutes)
	if err != nil {
		fmt.Println(err)
		return "", "", errors.New(err.Error())
	}

	expiredAt := now.Add(accessTokenExpiration)

	refreshExpiration, err := time.ParseDuration(s.cfg.JWT.RefreshExpirationHours)
	if err != nil {
		fmt.Println(err)
		return "", "", errors.New(err.Error())
	}

	refreshExpiratedAt := now.Add(refreshExpiration)

	claimsID := uuid.NewString()

	claims := auth.JWTClaims{
		UserID: userID,
		Email:  userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        claimsID,
			Issuer:    "laboratory-equipment-management-system",
			Subject:   userID,
			Audience:  []string{"laboratory-equipment-management-users"},
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.JWT.JwtSecret))
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	claims = auth.JWTClaims{
		UserID: userID,
		Email:  userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        claimsID,
			Issuer:    "refresh-token",
			Subject:   userID,
			Audience:  []string{"laboratory-equipment-management-users"},
			ExpiresAt: jwt.NewNumericDate(refreshExpiratedAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	// Refresh Token
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.cfg.JWT.RefreshSecret))
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
