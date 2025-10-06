package models

import (
	"time"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/google/uuid"
)

type Item struct {
	ItemID              uuid.UUID        `db:"item_id" gorm:"primaryKey;type:uuid"`
	ItemName            string           `db:"item_name"`
	ItemDescription     *string          `db:"item_description" gorm:"type:text"`
	ItemPictureURL      *string          `db:"item_picture_url" gorm:"type:text"`
	ItemStatus          enums.ItemStatus `db:"item_status"`
	ItemQuantity        int              `db:"item_quantity"`
	ItemCurrentQuantity int              `db:"item_current_quantity"`
	ItemCreatedAt       time.Time        `db:"item_created_at"`
	ItemUpdatedAt       time.Time        `db:"item_updated_at"`
}



type ItemParentChild struct {
	ParentID uuid.UUID `db:"parent_id" gorm:"type:uuid"`
	ChildID  uuid.UUID `db:"child_id" gorm:"type:uuid"`
}

func (Item) TableName() string { return "items" }
