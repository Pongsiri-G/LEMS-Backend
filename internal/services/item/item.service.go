package item

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	itemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	CreateItem(ctx context.Context, req *requests.CreateItemRequest) error
	GetBorrowItem(ctx context.Context, itemID string) (*responses.ItemResponse, error)
	GetAll(ctx context.Context) ([]responses.ItemResponse, error)
	GetMyBorrow(ctx context.Context, userID string) ([]responses.ItemResponse, error)
	GetChildItemByParentID(ctx context.Context, itemID string) ([]responses.ItemResponse, error)
}

type itemService struct {
	repo itemRepo.Repository
}

func NewItemService(repo itemRepo.Repository) Service {
	return &itemService{repo: repo}
}

func (i *itemService) GetChildItemByParentID(ctx context.Context, itemID string) ([]responses.ItemResponse, error) {
	itemIDUUID, err := uuid.Parse(itemID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return []responses.ItemResponse{}, ErrInvalidUUID
	}
	items, err := i.repo.GetChildItemByParentID(ctx, itemIDUUID)
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
		return &responses.ItemResponse{}, ErrInvalidUUID
	}

	item, err := i.repo.GetItemByID(ctx, itemIDUUID)

	if err != nil {
		return &responses.ItemResponse{}, err
	}

	response := responses.ItemResponse{
		ID:              item.ItemID,
		Name:            item.ItemName,
		Description:     item.ItemDescription,
		PictureURL:      item.ItemPictureURL,
		Status:          item.ItemStatus,
		Quantity:        item.ItemQuantity,
		CurrentQuantity: item.ItemCurrentQuantity,
		CreatedAt:       item.ItemCreatedAt,
		UpdatedAt:       item.ItemUpdatedAt,
	}
	return &response, nil
}

// CreateItem implements Service.
func (i *itemService) CreateItem(ctx context.Context, req *requests.CreateItemRequest) error {
	item := &models.Item{
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

func (i *itemService) GetAll(ctx context.Context) ([]responses.ItemResponse, error) {
	items, err := i.repo.GetAll(ctx)
	var response []responses.ItemResponse

	if err != nil {
		return nil, err
	}

	for _, i := range items {
		r := responses.ItemResponse{
			ID:          i.ItemID,
			Name:        i.ItemName,
			Description: i.ItemDescription,
			PictureURL:  i.ItemPictureURL,
			Status:      i.ItemStatus,
			Quantity:    i.ItemQuantity,
			CreatedAt:   i.ItemCreatedAt,
			UpdatedAt:   i.ItemUpdatedAt,
		}
		response = append(response, r)
	}

	return response, nil
}

func (i *itemService) GetMyBorrow(ctx context.Context, userID string) ([]responses.ItemResponse, error) {
	userUID, err := uuid.Parse(userID)

	if err != nil {
		log.Error().Err(err).Msg("Failed to parse UUID.")
		return nil, err
	}

	var items []models.Item
	items, err = i.repo.GetMyBorrow(ctx, userUID)

	if err != nil {
		return nil, err
	}

	var response []responses.ItemResponse
	for _, i := range items {
		r := responses.ItemResponse{
			ID:          i.ItemID,
			Name:        i.ItemName,
			Description: i.ItemDescription,
			PictureURL:  i.ItemPictureURL,
			Status:      i.ItemStatus,
			Quantity:    i.ItemQuantity,
			CreatedAt:   i.ItemCreatedAt,
			UpdatedAt:   i.ItemUpdatedAt,
		}
		response = append(response, r)
	}

	return response, nil
}
