package factory

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	ItemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	ItemSetRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item_set"
	TagRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/tag"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type ItemFactoryWithChildrenConcrete struct {
	request     *requests.CreateItemRequest
	itemRepo    ItemRepo.Repository
	itemSetRepo ItemSetRepo.Repository
	tagRepo     TagRepo.Repository
}

func NewItemFactoryWithChildrenConcrete(itemRepo ItemRepo.Repository, itemSetRepo ItemSetRepo.Repository, tagRepo TagRepo.Repository, request *requests.CreateItemRequest) ItemFactory {
	return &ItemFactoryWithChildrenConcrete{
		request:     request,
		itemRepo:    itemRepo,
		itemSetRepo: itemSetRepo,
		tagRepo:     tagRepo,
	}
}

// CreateItem implements ItemFactory.
func (i *ItemFactoryWithChildrenConcrete) CreateItem(ctx context.Context) (*models.Item, error) {

	var childrenIDs []uuid.UUID
	for _, itemID := range *i.request.Prerequisite {
		itemUUID, err := uuid.Parse(itemID)
		if err != nil {
			return nil, exceptions.ErrInvalidUUID
		}
		item, err := i.itemRepo.GetItemByID(ctx, itemUUID)
		if err != nil {
			return nil, err
		}

		if item == nil {
			return nil, exceptions.ErrItemNotFound
		}

		childrenIDs = append(childrenIDs, item.ItemID)
	}
	parentItem := models.Item{
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
		parentItem.ItemStatus = *i.request.Status
	}

	err := i.itemRepo.CreateItem(ctx, &parentItem)
	if err != nil {
		return nil, err
	}

	for _, childID := range childrenIDs {
		err := i.itemSetRepo.CreateItemSet(parentItem.ItemID, childID)
		// if any error occurs, rollback the parent item creation and any previously created item sets
		if err != nil {
			log.Error().Err(err).Msg("failed to create item set, rolling back parent item creation")
			for _, child := range childrenIDs {
				_ = i.itemSetRepo.DeleteItemSet(parentItem.ItemID, child)
			}
			_ = i.itemRepo.DeleteItem(ctx, parentItem.ItemID)
			return nil, err
		}
	}

	return &parentItem, nil

}
