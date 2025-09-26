package handlers

import authhd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/auth"

type Handlers struct {
	Auth *authhd.AuthHandler
}

func NewHandlers(auth *authhd.AuthHandler) *Handlers {
	return &Handlers{Auth: auth}
}
