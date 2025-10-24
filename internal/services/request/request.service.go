package request

import (
	ItemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	ItemRequested "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item_requested"
	RequestRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/request"
	RequestItem "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/request_item"
	"github.com/labstack/echo/v4"
)

type Service interface {
	CreateRequestService(c echo.Context) error
}

type service struct {
	requestRepo       RequestRepo.Repository
	requestItemRepo   RequestItem.Repository
	itemRequestedRepo ItemRequested.Repository
	itemRepo          ItemRepo.Repository
}

func NewRequestService(
	requestRepo RequestRepo.Repository,
	requestItemRepo RequestItem.Repository,
	itemRequestedRepo ItemRequested.Repository,
	itemRepo ItemRepo.Repository,
) Service {
	return &service{
		requestRepo:       requestRepo,
		requestItemRepo:   requestItemRepo,
		itemRequestedRepo: itemRequestedRepo,
		itemRepo:          itemRepo,
	}
}

// CreateRequestService implements Service.
func (s *service) CreateRequestService(c echo.Context) error {
	panic("unimplemented")
}
