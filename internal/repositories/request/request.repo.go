package request

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository interface {
	FindByID(ctx context.Context, requestID uuid.UUID) (*models.Request, error)
	CreateRequest(ctx context.Context, request *models.Request) error
	EditRequest(ctx context.Context, request *models.Request) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateRequest implements Repository.
func (r *repository) CreateRequest(ctx context.Context, request *models.Request) error {
	return r.db.WithContext(ctx).Create(request).Error
}

// FindByID implements Repository.
func (r *repository) FindByID(ctx context.Context, requestID uuid.UUID) (*models.Request, error) {
	var request models.Request
	if err := r.db.WithContext(ctx).First(&request, "request_id = ?", requestID).Error; err != nil {
		log.Error().Err(err).Msg("failed to find request by ID")
		return nil, err
	}
	return &request, nil
}

// EditRequest implements Repository.
func (r *repository) EditRequest(ctx context.Context, request *models.Request) error {
	return r.db.WithContext(ctx).Save(request).Error
}
