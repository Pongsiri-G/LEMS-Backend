package models

import (
	"time"

	"github.com/google/uuid"
)

type BorrowQueue struct {
	QueueID    uuid.UUID `db:"queue_id" gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `db:"user_id" gorm:"type:uuid;not null"`
	ItemID     uuid.UUID `db:"item_id" gorm:"type:uuid;not null"`
	IsBorrow   bool      `db:"is_borrow" gorm:"type:bool;default:false;not null"`
	CreatedAt  time.Time `db:"created_at" gorm:"not null"`
	BorrowedAt *time.Time `db:"borrowed_at" gorm:"default:null;null"`
}
