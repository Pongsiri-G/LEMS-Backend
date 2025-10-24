package itemrequested

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateItemRequested(ctx context.Context, item *models.ItemRequested) error
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
