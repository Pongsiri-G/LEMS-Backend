package itemrequested

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository interface {
	CreateItemRequested(ctx context.Context, item *models.ItemRequested) error
	FindByID(ctx context.Context, itemID uuid.UUID) (*models.ItemRequested, error)
	EditItemRequested(ctx context.Context, item *models.ItemRequested) error
}

type repository struct {
	db *gorm.DB
}

func NewItemRequestedRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateItemRequested implements Repository.
func (r *repository) CreateItemRequested(ctx context.Context, item *models.ItemRequested) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// FindByID implements Repository.
func (r *repository) FindByID(ctx context.Context, itemID uuid.UUID) (*models.ItemRequested, error) {
	var item models.ItemRequested
	if err := r.db.WithContext(ctx).First(&item, "id = ?", itemID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Msg("item requested not found for ID: " + itemID.String())
			return nil, nil
		}
		log.Error().Err(err).Msg("failed to find item requested by ID")
		return nil, err
	}
	return &item, nil
}

// EditItemRequested implements Repository.
func (r *repository) EditItemRequested(ctx context.Context, item *models.ItemRequested) error {
	return r.db.WithContext(ctx).Save(item).Error
}
