package borrowlog

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindBorrowLogByUserIDAndBorrowID(ctx context.Context, userID, borrowID uuid.UUID) (*models.BorrowLog, error)
	EditBorrowLog(ctx context.Context, borrowLog *models.BorrowLog) error
}

type repository struct {
	db *gorm.DB
}

func NewBorrowLogRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// FindBorrowLogByUserIDAndBorrowID implements Repository.
func (r *repository) FindBorrowLogByUserIDAndBorrowID(ctx context.Context, userID uuid.UUID, borrowID uuid.UUID) (*models.BorrowLog, error) {
	var borrowLogs []models.BorrowLog
	rows, err := r.db.Find(&models.BorrowLog{}, "user_id = ? AND borrow_id = ? AND borrow_status = ?", userID, borrowID, "BORROWED").Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	err = r.db.ScanRows(rows, &borrowLogs)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	if len(borrowLogs) > 1 {
		return nil, ErrMoreThanOneBorrowLog
	}

	return &borrowLogs[0], nil

}

// EditBorrowLog implements Repository.
func (r *repository) EditBorrowLog(ctx context.Context, borrowLog *models.BorrowLog) error {
	return r.db.WithContext(ctx).Save(borrowLog).Error
}
