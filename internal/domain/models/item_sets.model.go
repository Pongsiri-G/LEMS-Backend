package models

import (
	"github.com/google/uuid"
)

type ItemSets struct {
	ParentItemID uuid.UUID `db:"parent_item_id" gorm:"type:uuid;not null"`
	ChildItemID  uuid.UUID `db:"child_item_id" gorm:"type:uuid;not null"`
}
