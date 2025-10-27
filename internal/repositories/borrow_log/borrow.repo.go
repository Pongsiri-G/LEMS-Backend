package borrowlog

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindBorrowLogByUserIDAndBorrowID(ctx context.Context, userID, borrowID uuid.UUID) (*models.BorrowLog, error)
	EditBorrowLog(ctx context.Context, borrowLog *models.BorrowLog) error
	CreateBorrowLog(ctx context.Context, borrowLog models.BorrowLog) error
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]models.BorrowLog, error)
	CreateBorrowLogTx(ctx context.Context, borrowLog models.BorrowLog) error
}

type repository struct {
	db *gorm.DB
}

func NewBorrowLogRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateBorrowLog(ctx context.Context, borrowLog models.BorrowLog) error {
	return r.db.WithContext(ctx).Create(&borrowLog).Error
}

// FindBorrowLogByUserIDAndBorrowID implements Repository.
func (r *repository) FindBorrowLogByUserIDAndBorrowID(ctx context.Context, userID uuid.UUID, borrowID uuid.UUID) (*models.BorrowLog, error) {
	var borrowLog models.BorrowLog
	err := r.db.WithContext(ctx).Where("user_id = ? AND borrow_id = ? AND borrow_status = ?", userID, borrowID, "BORROWED").First(&borrowLog).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &borrowLog, nil
}

// EditBorrowLog implements Repository.
func (r *repository) EditBorrowLog(ctx context.Context, borrowLog *models.BorrowLog) error {
	return r.db.WithContext(ctx).Save(borrowLog).Error
}

// GetChildren implements Repository.
func (r *repository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]models.BorrowLog, error) {
	var borrowLogs []models.BorrowLog
	err := r.db.WithContext(ctx).Where("borrow_parent_id = ? AND borrow_status = ?", parentID, "BORROWED").Find(&borrowLogs).Error
	if err != nil {
		return nil, err
	}
	return borrowLogs, nil
}

// CreateBorrowLogTx implements Repository.
func (r *repository) CreateBorrowLogTx(ctx context.Context, borrowLog models.BorrowLog) error {
	tx := database.FromContext(ctx, r.db)
	return tx.Create(&borrowLog).Error
}
