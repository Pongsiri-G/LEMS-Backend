package jwt

import (
	"fmt"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/google/uuid"
)

type JWTService struct {
	cfg *configs.Config
}

func NewJWTService(cfg *configs.Config) *JWTService {
	return &JWTService{cfg: cfg}
}

// mock ไปก่อน
func (s *JWTService) Generate(userID uuid.UUID) (string, error) {
	return fmt.Sprintf("mock-token-%s", userID.String()), nil
}
