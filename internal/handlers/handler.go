package handlers

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/admin"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/auth"
	borrowhd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/borrow"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/borrowq"
	item "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/item"
	logHd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/log"
	miniohd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/minio"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/request"
	tag "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/tag"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/user"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/ws"
)

type Handlers struct {
	Admin       admin.AdminHandler
	Auth        auth.AuthHandler
	File        miniohd.FileHandler
	Borrow      borrowhd.BorrowHandler
	User        user.UserHandler
	Item        item.ItemHandler
	Tag         tag.TagHandler
	Request     request.RequestHandler
	Log         logHd.LogHandler
	BorrowQueue borrowq.BorrowQueueHandler
	WebSocket   ws.WsHandler
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
	log logHd.LogHandler,
	borrowQueue borrowq.BorrowQueueHandler,
	webSocket ws.WsHandler,
) *Handlers {
	return &Handlers{
		Admin:       admin,
		Auth:        auth,
		File:        file,
		Borrow:      borrow,
		User:        user,
		Item:        item,
		Tag:         tag,
		Request:     request,
		Log:         log,
		BorrowQueue: borrowQueue,
		WebSocket:   webSocket,
	}
}
