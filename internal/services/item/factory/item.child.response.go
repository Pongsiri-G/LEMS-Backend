package factory

import (
	"context"
	"sync"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	ItemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	ItemSetRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item_set"
	"github.com/rs/zerolog/log"
)

var childResponseLock = sync.Mutex{}
var childResponseFactory ItemResponseFactory = nil

type ItemResponseFactoryWithChildrenConcrete struct {
	itemRepo    ItemRepo.Repository
	itemSetRepo ItemSetRepo.Repository
}

func NewItemResponseFactoryWithChildrenConcrete(itemRepo ItemRepo.Repository, itemSetRepo ItemSetRepo.Repository) ItemResponseFactory {
	if childResponseFactory == nil {
		childResponseLock.Lock()
		defer childResponseLock.Unlock()
		if childResponseFactory == nil {
			childResponseFactory = &ItemResponseFactoryWithChildrenConcrete{
				itemRepo:    itemRepo,
				itemSetRepo: itemSetRepo,
			}
		}
	}
	return childResponseFactory
}

func (i *ItemResponseFactoryWithChildrenConcrete) ToResponse(ctx context.Context, item *models.Item, children *[]models.Item) (*responses.ItemResponse, error) {
	log.Info().Msg("Converting Item with Children to ItemResponse")
	var prerequisite []responses.ItemResponse
	visited := make(map[string]bool)

	// Mark parent as visited to avoid circular references
	visited[item.ItemID.String()] = true

	// Use DFS to collect all descendants recursively
	if err := i.collectAllDescendants(ctx, *children, &prerequisite, visited); err != nil {
		log.Error().Err(err).Msg("failed to collect all descendants")
		return nil, err
	}

	// Return the main item with all flattened prerequisites
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
		Prerequisites:   &prerequisite,
	}, nil
}

// collectAllDescendants recursively collects all child nodes using DFS
func (i *ItemResponseFactoryWithChildrenConcrete) collectAllDescendants(
	ctx context.Context,
	children []models.Item,
	result *[]responses.ItemResponse,
	visited map[string]bool,
) error {
	for _, child := range children {
		childID := child.ItemID.String()

		// Skip if already visited (handles cycles and duplicates)
		if visited[childID] {
			continue
		}

		// Mark as visited
		visited[childID] = true

		// Add current child to result
		*result = append(*result, responses.ItemResponse{
			ID:              child.ItemID,
			Name:            child.ItemName,
			Description:     child.ItemDescription,
			PictureURL:      child.ItemPictureURL,
			Status:          child.ItemStatus,
			Quantity:        child.ItemQuantity,
			CurrentQuantity: child.ItemCurrentQuantity,
			CreatedAt:       child.ItemCreatedAt,
			UpdatedAt:       child.ItemUpdatedAt,
		})

		// Fetch grandchildren from repository
		grandchildren, err := i.itemRepo.GetChildItemByParentID(ctx, child.ItemID)
		if err != nil {
			return err
		}

		// Recursively collect descendants of this child
		if len(grandchildren) > 0 {
			if err := i.collectAllDescendants(ctx, grandchildren, result, visited); err != nil {
				return err
			}
		}
	}

	return nil
}
