package models

import (
	"github.com/google/uuid"
)

type Tag struct {
	TagID    uuid.UUID `db:"tag_id" gorm:"type:uuid;primaryKey"`
	TagName  string    `db:"tag_name" gorm:"type:VARCHAR(20);not null"`
	TagColor string    `db:"tag_color" gorm:"type:VARCHAR(20);not null"`
}
