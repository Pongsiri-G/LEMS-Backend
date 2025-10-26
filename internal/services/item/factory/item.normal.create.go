package factory

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	ItemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	TagRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/tag"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type ItemFactoryConcrete struct {
	request  *requests.CreateItemRequest
	itemRepo ItemRepo.Repository
	tagRepo  TagRepo.Repository
}

func NewItemFactoryConcrete(itemRepo ItemRepo.Repository, tagRepo TagRepo.Repository, request *requests.CreateItemRequest) ItemFactory {
	return &ItemFactoryConcrete{
		request:  request,
		itemRepo: itemRepo,
		tagRepo:  tagRepo,
	}
}

// CreateItem implements ItemFactory.
func (i *ItemFactoryConcrete) CreateItem(ctx context.Context) (*models.Item, error) {
	item := &models.Item{
		ItemID:              uuid.New(),
		ItemName:            i.request.Name,
		ItemDescription:     i.request.Description,
		ItemPictureURL:      i.request.ImageURL,
		ItemStatus:          enums.ItemStatusAvailable,
		ItemQuantity:        i.request.Quantity,
		ItemCurrentQuantity: i.request.Quantity,
		ItemCreatedAt:       utils.BangkokNow(),
		ItemUpdatedAt:       utils.BangkokNow(),
	}

	if i.request.Status != nil {
		item.ItemStatus = *i.request.Status
	}

	err := i.itemRepo.CreateItem(ctx, item)
	if err != nil {
		log.Error().Err(err).Msg("failed to create item in database")
		return nil, err
	}

	return item, nil
}
