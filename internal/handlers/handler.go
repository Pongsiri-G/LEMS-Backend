package handlers

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/admin"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/auth"
	borrowhd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/borrow"
	item "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/item"
	miniohd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/minio"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/request"
	tag "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/tag"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/user"
)

type Handlers struct {
	Admin   admin.AdminHandler
	Auth    auth.AuthHandler
	File    miniohd.FileHandler
	Borrow  borrowhd.BorrowHandler
	User    user.UserHandler
	Item    item.ItemHandler
	Tag     tag.TagHandler
	Request request.RequestHandler
}

func NewHandlers(
	admin admin.AdminHandler,
	auth auth.AuthHandler,
	file miniohd.FileHandler,
	borrow borrowhd.BorrowHandler,
	user user.UserHandler,
	item item.ItemHandler,
	tag tag.TagHandler,
	request request.RequestHandler,
) *Handlers {
	return &Handlers{
		Admin:   admin,
		Auth:    auth,
		File:    file,
		Borrow:  borrow,
		User:    user,
		Item:    item,
		Tag:     tag,
		Request: request,
	}
}
