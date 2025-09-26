package jwt

import "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"

type JWTService struct {
	cfg *configs.Config
}

func NewJWTService(cfg *configs.Config) *JWTService {
	return &JWTService{cfg: cfg}
}
