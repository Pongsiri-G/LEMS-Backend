package request

import (
	"context"
	"time"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository interface {
	FindByID(ctx context.Context, requestID uuid.UUID) (*models.Request, error)
	CreateRequest(ctx context.Context, request *models.Request) error
	EditRequest(ctx context.Context, request *models.Request) error

	GetRequests(ctx context.Context, requestType *enums.RequestType, requestStatus *enums.RequestStatus) ([]models.Request, error)
	GetRequestsByUserID(ctx context.Context, userID uuid.UUID, requestType *enums.RequestType, requestStatus *enums.RequestStatus) ([]models.Request, error)
	ChangeRequestStatus(ctx context.Context, requestID uuid.UUID, status enums.RequestStatus) error
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
	request.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(request).Error
}

// GetRequests implements Repository.
func (r *repository) GetRequests(ctx context.Context, requestType *enums.RequestType, requestStatus *enums.RequestStatus) ([]models.Request, error) {

	var requests []models.Request
	database := r.db.WithContext(ctx)

	if requestType != nil {
		database = database.Where("request_type = ?", *requestType)
	}

	if requestStatus != nil {
		database = database.Where("request_status = ?", *requestStatus)
	}

	if err := database.Find(&requests).Order("updated_at desc").Error; err != nil {
		log.Error().Err(err).Msg("failed to get requests by type")
		return nil, err
	}

	return requests, nil
}

// GetRequestsByUserID implements Repository.
func (r *repository) GetRequestsByUserID(ctx context.Context, userID uuid.UUID, requestType *enums.RequestType, requestStatus *enums.RequestStatus) ([]models.Request, error) {
	var requests []models.Request
	database := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if requestType != nil {
		database = database.Where("request_type = ?", *requestType)
	}

	if requestStatus != nil {
		database = database.Where("request_status = ?", *requestStatus)
	}

	if err := database.Find(&requests).Order("updated_at desc").Error; err != nil {
		log.Error().Err(err).Msg("failed to get requests by user ID")
		return nil, err
	}

	return requests, nil
}

// ChangeRequestStatus implements Repository.
func (r *repository) ChangeRequestStatus(ctx context.Context, requestID uuid.UUID, status enums.RequestStatus) error {
	return r.db.WithContext(ctx).Model(&models.Request{}).Where("request_id = ?", requestID).Update("request_status", status).Error
}
