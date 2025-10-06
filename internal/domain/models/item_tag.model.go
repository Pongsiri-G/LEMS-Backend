package models

import (
	"github.com/google/uuid"
)

type ItemTag struct {
	ItemID       uuid.UUID          `db:"item_id" gorm:"type:uuid;not null"`
	TagID       uuid.UUID          `db:"tag_id" gorm:"type:uuid;not null"`
}
