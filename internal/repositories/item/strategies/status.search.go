package repository

import (
	"strings"

	"gorm.io/gorm"
)

type StatusSearch struct {
    Status string
}

func (s StatusSearch) Apply(db *gorm.DB) *gorm.DB {
    if s.Status == "" {
        return db
    }
    return db.Where("LOWER(item_status) = ?", strings.ToLower(s.Status))
}