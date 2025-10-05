package handlers

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/auth"
	borrowhd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/borrow"
	item "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/item"
	miniohd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/minio"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/user"
)
type Handlers struct {
	Auth   auth.AuthHandler
	File   miniohd.FileHandler
	Borrow borrowhd.BorrowHandler
	User user.UserHandler
	Item   item.ItemHandler
}

func NewHandlers(
	auth auth.AuthHandler,
	file miniohd.FileHandler,
	borrow borrowhd.BorrowHandler,
	user user.UserHandler,
	item item.ItemHandler,
) *Handlers {
	return &Handlers{
		Auth:   auth,
		File:   file,
		Borrow: borrow,
		User: user,
		Item:   item,
	}
}
