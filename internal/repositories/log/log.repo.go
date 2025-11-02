package log

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, log *models.Log) error
	List(ctx context.Context) ([]models.Log, error)
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{db: db}
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
