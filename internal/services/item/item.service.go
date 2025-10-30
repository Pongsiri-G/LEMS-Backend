package item

import (
	"context"
	"strings"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	ItemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	repository "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item/strategies"
	ItemSetRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item_set"
	TagRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/tag"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/item/factory"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils/itemutil"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	CreateItem(ctx context.Context, req *requests.CreateItemRequest) error
	UpdateItem(ctx context.Context, req *requests.EditItemRequest) error

	GetBorrowItem(ctx context.Context, itemID string) (*responses.ItemResponse, error)
	GetAll(ctx context.Context) ([]responses.ItemResponse, error)
	GetMyBorrow(ctx context.Context, userID string) ([]responses.ItemResponseBorrow, error)
	GetChildItemByParentID(ctx context.Context, itemID string) ([]responses.ItemResponse, error)
	SearchItems(ctx context.Context, strategies ItemRepo.SearchStrategyMap) ([]responses.ItemResponse, error)
	DeleteItem(ctx context.Context, itemID string) error

	AssignChild(ctx context.Context, parentID, childID string) error
	RemoveChild(ctx context.Context, parentID, childID string) error
}

type itemService struct {
	itemRepo    ItemRepo.Repository
	itemSetRepo ItemSetRepo.Repository
	tagRepo     TagRepo.Repository
}

func NewItemService(itemRepo ItemRepo.Repository, itemSetRepo ItemSetRepo.Repository, tagRepo TagRepo.Repository) Service {
	return &itemService{itemRepo: itemRepo, itemSetRepo: itemSetRepo, tagRepo: tagRepo}
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

	if req.Tags != nil && len(*req.Tags) > 0 {
		for _, tagIDStr := range *req.Tags {
			tagID, err := uuid.Parse(tagIDStr)
			if err != nil {
				return exceptions.ErrInvalidUUID
			}
			tagModel, err := i.tagRepo.GetTagByID(ctx, tagID)
			if err != nil {
				log.Error().Err(err).Msg("failed to get tag by id")
				return err
			}
			if tagModel == nil {
				log.Error().Msgf("tag not found: %s", tagIDStr)
				return exceptions.ErrTagNotFound
			}
		}
	}

	if req.Prerequisite != nil && len(*req.Prerequisite) > 0 {
		itemFactory = factory.NewItemFactoryWithChildrenConcrete(i.itemRepo, i.itemSetRepo, i.tagRepo, req)
	} else {
		itemFactory = factory.NewItemFactoryConcrete(i.itemRepo, i.tagRepo, req)
	}
	item, err := itemFactory.CreateItem(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to create item")
		return err
	}

	if req.Tags != nil && len(*req.Tags) > 0 {
		for _, tagIDStr := range *req.Tags {
			tagID, _ := uuid.Parse(tagIDStr)
			err := i.tagRepo.AssignTagToItem(ctx, item.ItemID, tagID)
			if err != nil {
				log.Error().Err(err).Msg("failed to assign tag to item")
				return err
			}
		}
	}

	return nil

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

func (i *itemService) SearchItems(ctx context.Context, strategiesMap ItemRepo.SearchStrategyMap) ([]responses.ItemResponse, error) {
	var tagsCleaned []string
	unique := map[string]struct{}{}
	for _, tag := range strategiesMap.Tags {
		for _, t := range strings.Split(tag, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				unique[t] = struct{}{}
			}
		}
	}

	for tag := range unique {
		tagsCleaned = append(tagsCleaned, tag)
	}

	strategies := []ItemRepo.SearchStrategy{
		repository.NameSearch{Query: strategiesMap.Name},
		repository.TagSearch{Tags: tagsCleaned},
		repository.StatusSearch{Status: strategiesMap.Status},
		repository.UserSearch{Query: strategiesMap.User},
	}

	log.Debug().Msgf("query := name: %s, tags: %s, status: %s", strategiesMap.Name, strings.Join(tagsCleaned, ","), strategiesMap.Status)
	items, err := i.itemRepo.SearchItems(ctx, strategies)

	if err != nil {
		return nil, err
	}

	return itemutil.ToResponses(items), nil

}

// DeleteItem implements Service.
func (i *itemService) DeleteItem(ctx context.Context, itemID string) error {
	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return exceptions.ErrInvalidUUID
	}

	item, err := i.itemRepo.GetItemByID(ctx, itemUUID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get item by ID")
		return err
	}

	if item == nil {
		log.Error().Msg("item not found for ID: " + itemID)
		return exceptions.ErrItemNotFound
	}

	err = i.itemRepo.DeleteItem(ctx, itemUUID)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete item")
		return err
	}

	return nil
}

// AssignChild implements Service.
func (i *itemService) AssignChild(ctx context.Context, parentID string, childID string) error {
	parentUUID, err := uuid.Parse(parentID)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse parent ID")
		return exceptions.ErrInvalidUUID
	}

	childUUID, err := uuid.Parse(childID)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse child ID")
		return exceptions.ErrInvalidUUID
	}

	itemSet, err := i.itemSetRepo.FindItemSetByParentAndChildID(ctx, parentUUID, childUUID)
	if err != nil {
		log.Error().Err(err).Msg("failed to find item set by parent and child ID")
		return err
	}

	if itemSet != nil {
		log.Info().Msg("item set already exists")
		return exceptions.ErrItemSetAlreadyExists
	}

	err = i.itemSetRepo.CreateItemSet(ctx, parentUUID, childUUID)
	if err != nil {
		log.Error().Err(err).Msg("failed to assign child item")
		return err
	}

	return nil
}

// RemoveChild implements Service.
func (i *itemService) RemoveChild(ctx context.Context, parentID string, childID string) error {
	parentUUID, err := uuid.Parse(parentID)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse parent ID")
		return exceptions.ErrInvalidUUID
	}

	childUUID, err := uuid.Parse(childID)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse child ID")
		return exceptions.ErrInvalidUUID
	}

	err = i.itemSetRepo.DeleteItemSet(ctx, parentUUID, childUUID)
	if err != nil {
		log.Error().Err(err).Msg("failed to remove child item")
		return err
	}

	return nil
}

// UpdateItem implements Service.
func (i *itemService) UpdateItem(ctx context.Context, req *requests.EditItemRequest) error {
	itemID, err := uuid.Parse(req.ItemID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return exceptions.ErrInvalidUUID
	}

	item, err := i.itemRepo.GetItemByID(ctx, itemID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get item by ID")
		return err
	}

	if item == nil {
		log.Error().Msg("item not found for ID: " + req.ItemID)
		return exceptions.ErrItemNotFound
	}

	if req.Name != nil {
		item.ItemName = *req.Name
	}
	if req.Description != nil {
		item.ItemDescription = req.Description
	}
	if req.ImageURL != nil {
		item.ItemPictureURL = req.ImageURL
	}
	if req.Quantity != nil {
		item.ItemQuantity = *req.Quantity
	}
	if req.Status != nil {
		item.ItemStatus = *req.Status
	}

	err = i.itemRepo.UpdateItem(ctx, item)
	if err != nil {
		log.Error().Err(err).Msg("failed to update item")
		return err
	}

	return nil
}
