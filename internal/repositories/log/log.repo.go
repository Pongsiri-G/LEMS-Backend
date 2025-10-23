package log

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, log *models.Log) error
	CreateBorrowLog(ctx context.Context, userID, itemID uuid.UUID) error
	CreateReturnLog(ctx context.Context, userID, itemID uuid.UUID) error
	List(ctx context.Context) ([]models.Log, error)
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{db: db}
}

// CreateBorrowLog implements Repository.
func (r *RepositoryImpl) CreateBorrowLog(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error {
	panic("unimplemented")
}

// CreateReturnLog implements Repository.
func (r *RepositoryImpl) CreateReturnLog(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error {
	panic("unimplemented")
}

func (r *RepositoryImpl) Create(ctx context.Context, log *models.Log) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *RepositoryImpl) List(ctx context.Context) ([]models.Log, error) {
	var logs []models.Log
	q := r.db.WithContext(ctx).Find(&logs)
	if q.Error != nil {
		return nil, q.Error
	}
	return logs, nil
}
