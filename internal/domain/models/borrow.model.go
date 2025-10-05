package models

import (
	"time"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/google/uuid"
)

type BorrowLog struct {
	BorrowID     uuid.UUID          `db:"borrow_id" gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID          `db:"user_id" gorm:"type:uuid;not null"`
	ItemID       uuid.UUID          `db:"item_id" gorm:"type:uuid;not null"`
	BorrowStatus enums.BorrowStatus `db:"borrow_status" gorm:"type:VARCHAR(20);not null"`
	BorrowDate   time.Time          `db:"borrow_date" gorm:"not null"`
	ReturnDate   *time.Time         `db:"return_date" gorm:"default:null"`
	CreatedAt    time.Time          `db:"created_at" gorm:"not null"`
	UpdatedAt    time.Time          `db:"updated_at" gorm:"not null"`
}
