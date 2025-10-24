package requestitem

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	AssignRequestItemToRequest(ctx context.Context, requestID, itemID uuid.UUID, quantity int) error
	DeleteRequestItemFromRequest(ctx context.Context, requestItemID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRequestItemRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// AssignRequestItemToRequest implements Repository.
func (r *repository) AssignRequestItemToRequest(ctx context.Context, requestID uuid.UUID, itemID uuid.UUID, quantity int) error {
	requestItem := models.RequestItem{
		RequestItemID: uuid.New(),
		RequestID:     requestID,
		ItemID:        itemID,
		Quantity:      quantity,
	}

	return r.db.WithContext(ctx).Create(&requestItem).Error
}

// DeleteRequestItemFromRequest implements Repository.
func (r *repository) DeleteRequestItemFromRequest(ctx context.Context, requestItemID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.RequestItem{}, "request_item_id = ?", requestItemID).Error
}
