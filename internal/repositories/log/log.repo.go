package log

import (
	"context"
	"encoding/json"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
	log.Info().Msg("Creating borrow log")
	jsonMap := map[string]uuid.UUID{
		"user_id": userID,
		"item_id": itemID,
	}
	logBytes, err := json.Marshal(jsonMap)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal log message to JSON")
		return err
	}
	logMessage := string(logBytes)
	logEntry := &models.Log{
		LogID:      uuid.New(),
		UserID:     userID,
		LogType:    enums.LogTypeBorrow,
		LogMessage: &logMessage,
	}
	return r.db.WithContext(ctx).Create(logEntry).Error

}

// CreateReturnLog implements Repository.
func (r *RepositoryImpl) CreateReturnLog(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error {
	log.Info().Msg("Creating return log")
	jsonMap := map[string]uuid.UUID{
		"user_id": userID,
		"item_id": itemID,
	}

	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal log message to JSON")
		return err
	}
	logMessage := string(jsonBytes)
	logEntry := &models.Log{
		LogID:      uuid.New(),
		UserID:     userID,
		LogType:    enums.LogTypeReturn,
		LogMessage: &logMessage,
	}
	return r.db.WithContext(ctx).Create(logEntry).Error
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
