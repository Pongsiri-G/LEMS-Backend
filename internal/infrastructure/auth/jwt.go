package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}
