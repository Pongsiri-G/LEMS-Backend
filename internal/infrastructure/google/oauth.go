package google

import "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"

type GoogleOAuthClient struct {
	cfg *configs.Config
}

func NewGoogleOAuthClient(cfg *configs.Config) *GoogleOAuthClient {
	return &GoogleOAuthClient{cfg: cfg}
}
