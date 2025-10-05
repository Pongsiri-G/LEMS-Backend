package handlers

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/auth"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/user"
)
type Handlers struct {
	Auth   authhd.AuthHandler
	File   miniohd.FileHandler
	Borrow borrowhd.BorrowHandler
	User user.UserHandler
}

func NewHandlers(
	auth authhd.AuthHandler,
	file miniohd.FileHandler,
	borrow borrowhd.BorrowHandler,
	user user.UserHandler,
) *Handlers {
	return &Handlers{
		Auth:   auth,
		File:   file,
		Borrow: borrow,
		User: user,
	}
}
