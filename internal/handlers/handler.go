package handlers

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/user"
)

type Handlers struct {
	Auth auth.AuthHandler
	User user.UserHandler
}

func NewHandlers(
	auth auth.AuthHandler,
	user user.UserHandler,
) *Handlers {
	return &Handlers{
		Auth: auth,
		User: user,
	}
}
