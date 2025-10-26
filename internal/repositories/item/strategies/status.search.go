package repository

import "gorm.io/gorm"

type StatusSearch struct {
    Status string
}

func (s StatusSearch) Apply(db *gorm.DB) *gorm.DB {
    if s.Status == "" {
        return db
    }
    return db.Where("item_status = ?", s.Status)
}