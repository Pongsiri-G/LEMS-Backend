package item

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	itemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/item/strategy"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils/itemutil"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	CreateItem(ctx context.Context, req *requests.CreateItemRequest) error
	GetBorrowItem(ctx context.Context, itemID string) (*responses.ItemResponse, error)
	GetAll(ctx context.Context) ([]responses.ItemResponse, error)
	GetMyBorrow(ctx context.Context, userID string) ([]responses.ItemResponse, error)
	GetChildItemByParentID(ctx context.Context, itemID string) ([]responses.ItemResponse, error)
	GetFiltered(ctx context.Context, strategy string, query []string) ([]responses.ItemResponse, error)
}

type itemService struct {
	repo itemRepo.Repository
	f map[string]strategy.FilterStrategy
}

func NewItemService(repo itemRepo.Repository) Service {
	return &itemService{repo: repo, f: strategy.NewFilterMap(nil)}
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

	response := itemutil.ToResponse(*item)
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

	if err != nil {
		return nil, err
	}

	response := itemutil.ToResponses(items)	

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

	response := itemutil.ToResponses(items)

	return response, nil
}

func (i *itemService) GetFiltered(ctx context.Context, strat string, query []string) ([]responses.ItemResponse, error) {
	i.f = strategy.NewFilterMap(query)
	
	filter := i.f[strat]

	if filter == nil {
		return nil, exceptions.ErrNoSuchStrategy
	}

	filter.InitFilter(i.repo)

	items, err := filter.Filter(ctx)

	if err != nil {
		return nil, err
	}
	response := itemutil.ToResponses(items)

	return response, nil
}
