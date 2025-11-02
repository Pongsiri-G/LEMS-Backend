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
	CreateLoginLog(ctx context.Context, userID uuid.UUID, message string) error
	CreateRegisterLog(ctx context.Context, userID uuid.UUID) error
	CreateAdminActionLog(ctx context.Context, adminID uuid.UUID, actionType enums.LogType, targetUserID uuid.UUID) error
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

// CreateLoginLog implements Repository.
func (r *RepositoryImpl) CreateLoginLog(ctx context.Context, userID uuid.UUID, message string) error {
	log.Info().Msg("Creating login log")
	logEntry := &models.Log{
		LogID:      uuid.New(),
		UserID:     userID,
		LogType:    enums.LogTypeLogin,
		LogMessage: &message,
	}
	return r.db.WithContext(ctx).Create(logEntry).Error
}

// CreateRegisterLog implements Repository.
func (r *RepositoryImpl) CreateRegisterLog(ctx context.Context, userID uuid.UUID) error {
	log.Info().Msg("Creating register log")
	message := "User registered successfully"
	logEntry := &models.Log{
		LogID:      uuid.New(),
		UserID:     userID,
		LogType:    enums.LogTypeRegister,
		LogMessage: &message,
	}
	return r.db.WithContext(ctx).Create(logEntry).Error
}

// CreateAdminActionLog implements Repository.
func (r *RepositoryImpl) CreateAdminActionLog(ctx context.Context, adminID uuid.UUID, actionType enums.LogType, targetUserID uuid.UUID) error {
	log.Info().Msgf("Creating admin action log: %s", actionType)
	message := targetUserID.String()
	logEntry := &models.Log{
		LogID:      uuid.New(),
		UserID:     adminID,
		LogType:    actionType,
		LogMessage: &message,
	}
	return r.db.WithContext(ctx).Create(logEntry).Error
}

func (r *RepositoryImpl) List(ctx context.Context) ([]models.Log, error) {
	var logs []models.Log
	q := r.db.WithContext(ctx).Find(&logs)
	if q.Error != nil {
		return nil, q.Error
	}
	return logs, nil
}
