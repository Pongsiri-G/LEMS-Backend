package factory

import (
	"context"
	"sync"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	ItemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
)

var normalFactoryLock = sync.Mutex{}
var normalFactory ItemFactory = nil

type ItemFactoryConcrete struct {
	request  *requests.CreateItemRequest
	itemRepo ItemRepo.Repository
}

func NewItemFactoryConcrete(itemRepo ItemRepo.Repository, request *requests.CreateItemRequest) ItemFactory {
	if normalFactory == nil {
		normalFactoryLock.Lock()
		defer normalFactoryLock.Unlock()
		if normalFactory == nil {
			normalFactory = &ItemFactoryConcrete{
				itemRepo: itemRepo,
				request:  request,
			}
		}
	}
	return normalFactory
}

// CreateItem implements ItemFactory.
func (i *ItemFactoryConcrete) CreateItem(ctx context.Context) error {
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

	return i.itemRepo.CreateItem(ctx, item)
}
