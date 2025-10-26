package item

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	ItemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	ItemSetRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item_set"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/item/factory"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/item/strategy"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils/itemutil"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	CreateItem(ctx context.Context, req *requests.CreateItemRequest) error
	GetBorrowItem(ctx context.Context, itemID string) (*responses.ItemResponse, error)
	GetAll(ctx context.Context) ([]responses.ItemResponse, error)
	GetMyBorrow(ctx context.Context, userID string) ([]responses.ItemResponseBorrow, error)
	GetChildItemByParentID(ctx context.Context, itemID string) ([]responses.ItemResponse, error)
	GetFiltered(ctx context.Context, strategy string, query []string) ([]responses.ItemResponse, error)
}

type itemService struct {
	itemRepo    ItemRepo.Repository
	itemSetRepo ItemSetRepo.Repository
	f           map[string]strategy.FilterStrategy
}

// GetFiltered implements Service.
func (i *itemService) GetFiltered(ctx context.Context, strategy string, query []string) ([]responses.ItemResponse, error) {
	panic("unimplemented")
}

func NewItemService(itemRepo ItemRepo.Repository, itemSetRepo ItemSetRepo.Repository) Service {
	return &itemService{itemRepo: itemRepo, itemSetRepo: itemSetRepo, f: strategy.NewFilterMap(nil)}
}

func (i *itemService) GetChildItemByParentID(ctx context.Context, itemID string) ([]responses.ItemResponse, error) {
	itemIDUUID, err := uuid.Parse(itemID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return []responses.ItemResponse{}, exceptions.ErrInvalidUUID
	}
	items, err := i.itemRepo.GetChildItemByParentID(ctx, itemIDUUID)
	response := make([]responses.ItemResponse, 0)

	if err != nil {
		return nil, err
	}

	for _, i := range items {
		r := responses.ItemResponse{
			ID:              i.ItemID,
			Name:            i.ItemName,
			Description:     i.ItemDescription,
			PictureURL:      i.ItemPictureURL,
			Status:          i.ItemStatus,
			Quantity:        i.ItemQuantity,
			CurrentQuantity: i.ItemCurrentQuantity,
			CreatedAt:       i.ItemCreatedAt,
			UpdatedAt:       i.ItemUpdatedAt,
		}
		response = append(response, r)
	}

	return response, nil
}

func (i *itemService) GetBorrowItem(ctx context.Context, itemID string) (*responses.ItemResponse, error) {
	itemIDUUID, err := uuid.Parse(itemID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return nil, exceptions.ErrInvalidUUID
	}

	item, err := i.itemRepo.GetItemByID(ctx, itemIDUUID)
	if item == nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	response := itemutil.ToResponse(*item)
	return &response, nil
}

// CreateItem implements Service.
func (i *itemService) CreateItem(ctx context.Context, req *requests.CreateItemRequest) error {
	var itemFactory factory.ItemFactory
	if req.Prerequisite != nil && len(*req.Prerequisite) > 0 {
		itemFactory = factory.NewItemFactoryWithChildrenConcrete(i.itemRepo, i.itemSetRepo, req)
	} else {
		itemFactory = factory.NewItemFactoryConcrete(i.itemRepo, req)
	}
	return itemFactory.CreateItem(ctx)
}

func (i *itemService) GetAll(ctx context.Context) ([]responses.ItemResponse, error) {
	childFactory := factory.NewItemResponseFactoryWithChildrenConcrete(i.itemRepo, i.itemSetRepo)
	normalFactory := factory.NewItemResponseFactoryConcrete(i.itemRepo)
	items, err := i.itemRepo.GetAll(ctx)

	if err != nil {
		return nil, err
	}
	var res []responses.ItemResponse
	for _, item := range items {
		children, err := i.itemRepo.GetChildItemByParentID(ctx, item.ItemID)
		if err != nil {
			return nil, err
		}

		var response *responses.ItemResponse
		if len(children) > 0 {
			response, err = childFactory.ToResponse(ctx, &item, &children)
			if err != nil {
				return nil, err
			}
		} else {
			response, err = normalFactory.ToResponse(ctx, &item, nil)
			if err != nil {
				return nil, err
			}
		}

		res = append(res, *response)

	}

	return res, nil
}

func (i *itemService) GetMyBorrow(ctx context.Context, userID string) ([]responses.ItemResponseBorrow, error) {
	userUID, err := uuid.Parse(userID)

	if err != nil {
		log.Error().Err(err).Msg("Failed to parse UUID.")
		return nil, err
	}

	var items []models.ItemBorrow
	items, err = i.itemRepo.GetMyBorrow(ctx, userUID)

	if err != nil {
		return nil, err
	}

	var response []responses.ItemResponseBorrow
	for _, i := range items {
		r := responses.ItemResponseBorrow{
			ID:          i.ItemID,
			Name:        i.ItemName,
			Description: i.ItemDescription,
			PictureURL:  i.ItemPictureURL,
			Status:      i.ItemStatus,
			Quantity:    i.ItemQuantity,
			CreatedAt:   i.ItemCreatedAt,
			UpdatedAt:   i.ItemUpdatedAt,
			BorrowID:    i.BorrowID,
		}
		response = append(response, r)
	}

	return response, nil
}
