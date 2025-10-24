package models

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/google/uuid"
)

type RequestItem struct {
	RequestItemID     uuid.UUID               `db:"request_item_id" gorm:"primaryKey;type:uuid"`
	RequestID         uuid.UUID               `db:"request_id" gorm:"type:uuid;not null"`
	ItemID            uuid.UUID               `db:"item_id" gorm:"type:uuid;not null"`
	Quantity          int                     `db:"quantity" gorm:"not null"`
	RequestItemStatus enums.RequestItemStatus `db:"request_item_status" gorm:"not null"`
}
