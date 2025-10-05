package item

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	itemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
)

type Service interface {
	CreateItem(ctx context.Context, req *requests.CreateItemRequest) error
}

type itemService struct {
	repo itemRepo.Repository
}

func NewItemService(repo itemRepo.Repository) Service {
	return &itemService{repo: repo}
}

// CreateItem implements Service.
func (i *itemService) CreateItem(ctx context.Context, req *requests.CreateItemRequest) error {
	item := &models.Items{
		ItemID:          uuid.New(),
		ItemName:        req.Name,
		ItemDescription: req.Description,
		ItemPictureURL:  req.ImageURL,
		ItemQuantity:    req.Quantity,
		ItemStatus:      enums.ItemStatusAvailable,
		ItemCreatedAt:   utils.BangkokNow(),
		ItemUpdatedAt:   utils.BangkokNow(),
	}

	return i.repo.CreateItem(ctx, item)

}
