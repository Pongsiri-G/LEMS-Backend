package models

import (
	"time"

	"github.com/google/uuid"
)

type Items struct {
	ItemID          uuid.UUID `db:"item_id" gorm:"primaryKey"`
	ItemName        string    `db:"item_name"`
	ItemDescription string    `db:"item_description"`
	ItemPictureURL  string    `db:"item_picture_url"`
	ItemStatus      string    `db:"item_status"`
	ItemQuantity    int       `db:"item_quantity"`
	ItemCreatedAt   time.Time `db:"item_created_at"`
	ItemUpdatedAt   time.Time `db:"item_updated_at"`
}

func (Items) TableName() string { return "items" }
