package handlers

import handlers "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/auth"

type Handlers struct {
	Auth handlers.AuthHandler
}

func NewHandlers(
	authHandler handlers.AuthHandler,
) *Handlers {
	return &Handlers{
		Auth: authHandler,
	}
}
