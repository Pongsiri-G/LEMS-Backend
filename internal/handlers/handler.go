package handlers

import (
	authhd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/auth"
	miniohd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/minio"
)

type Handlers struct {
	Auth authhd.AuthHandler
	File miniohd.FileHandler
}

func NewHandlers(
	auth authhd.AuthHandler,
	file miniohd.FileHandler,
) *Handlers {
	return &Handlers{
		Auth: auth,
		File: file,
	}
}
