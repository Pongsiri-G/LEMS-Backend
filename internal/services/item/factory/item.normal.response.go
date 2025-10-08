package factory

import (
	"context"
	"sync"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	ItemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	"github.com/rs/zerolog/log"
)

var normalResponseLock = sync.Mutex{}
var normalResponseFactory ItemResponseFactory = nil

type ItemResponseFactoryConcrete struct {
	itemRepo ItemRepo.Repository
}

func NewItemResponseFactoryConcrete(itemRepo ItemRepo.Repository) ItemResponseFactory {
	if normalResponseFactory == nil {
		normalResponseLock.Lock()
		defer normalResponseLock.Unlock()
		if normalResponseFactory == nil {
			normalResponseFactory = &ItemResponseFactoryConcrete{
				itemRepo: itemRepo,
			}
		}
	}
	return normalResponseFactory
}

// ToResponse implements ItemResponseFactory.
func (i *ItemResponseFactoryConcrete) ToResponse(ctx context.Context, item *models.Item, children *[]models.Item) (*responses.ItemResponse, error) {
	log.Info().Msg("Converting Item to ItemResponse")
	return &responses.ItemResponse{
		ID:              item.ItemID,
		Name:            item.ItemName,
		Description:     item.ItemDescription,
		PictureURL:      item.ItemPictureURL,
		Status:          item.ItemStatus,
		Quantity:        item.ItemQuantity,
		CurrentQuantity: item.ItemCurrentQuantity,
		CreatedAt:       item.ItemCreatedAt,
		UpdatedAt:       item.ItemUpdatedAt,
	}, nil
}
