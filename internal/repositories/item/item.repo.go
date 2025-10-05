package item

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateItem(ctx context.Context, item *models.Items) error
}

type repository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateItem implements Repository.
func (r *repository) CreateItem(ctx context.Context, item *models.Items) error {
	return r.db.WithContext(ctx).Create(item).Error
}
