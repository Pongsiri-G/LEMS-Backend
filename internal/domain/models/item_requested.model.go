package models

import "github.com/google/uuid"

type ItemRequested struct {
	ID          uuid.UUID `db:"item_requested_id" gorm:"primaryKey;type:uuid"`
	Name        string    `db:"item_requested_name" gorm:"not null"`
	Description string    `db:"item_requested_description" gorm:"not null"`
	Type        string    `db:"item_requested_type" gorm:"not null"`
	UserID      uuid.UUID `db:"user_id" gorm:"type:uuid;not null"`
	Quantity    int       `db:"quantity" gorm:"not null"`
	Price       float64   `db:"price" gorm:"not null"`
}

func (ItemRequested) TableName() string { return "items_requested" }
