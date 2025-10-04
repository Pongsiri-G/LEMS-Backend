package handlers

import (
	authhd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/auth"
	borrowhd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/borrow"
	item "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/item"
	miniohd "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/handlers/minio"
)

type Handlers struct {
	Auth   authhd.AuthHandler
	File   miniohd.FileHandler
	Borrow borrowhd.BorrowHandler
	Item   item.ItemHandler
}

func NewHandlers(
	auth authhd.AuthHandler,
	file miniohd.FileHandler,
	borrow borrowhd.BorrowHandler,
	item item.ItemHandler,
) *Handlers {
	return &Handlers{
		Auth:   auth,
		File:   file,
		Borrow: borrow,
		Item:   item,
	}
}
