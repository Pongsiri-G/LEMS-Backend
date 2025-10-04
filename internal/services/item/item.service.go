package item

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	itemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	CreateItem(ctx context.Context, req *requests.CreateItemRequest) error
	GetBorrowItem(ctx context.Context, itemID string) (*models.Items, error)
}

type itemService struct {
	repo itemRepo.Repository
}

func NewItemService(repo itemRepo.Repository) Service {
	return &itemService{repo: repo}
}

func (i *itemService) GetBorrowItem(ctx context.Context, itemID string) (*models.Items, error) {
	itemIDUUID, err := uuid.Parse(itemID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return &models.Items{}, ErrInvalidUUID
	}

	item, err := i.repo.GetItemByID(ctx, itemIDUUID)

	if err != nil {
		return &models.Items{}, err
	}

	return item, nil
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
