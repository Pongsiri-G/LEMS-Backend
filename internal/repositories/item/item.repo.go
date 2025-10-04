package item

import (
	"context"
	"fmt"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateItem(ctx context.Context, item *models.Items) error
	GetItemByID(ctx context.Context, itemID uuid.UUID) (*models.Items, error)
}

type repository struct {
	db *gorm.DB
}

func (r *repository) GetItemByID(ctx context.Context, itemID uuid.UUID) (*models.Items, error) {
	var item models.Items
	err := r.db.WithContext(ctx).Where("item_id = ?", itemID).First(&item).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println(("ERR Record not found"))
			return nil, err
		}
		return nil, err
	}

	return &item, nil

}

func NewItemRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateItem implements Repository.
func (r *repository) CreateItem(ctx context.Context, item *models.Items) error {
	return r.db.WithContext(ctx).Create(item).Error
}
