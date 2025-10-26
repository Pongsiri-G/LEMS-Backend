package repository

import "gorm.io/gorm"

type NameSearch struct {
	Query string
}

func (s NameSearch) Apply(db *gorm.DB) *gorm.DB {
	if s.Query == "" {
		return db
	}
	return db.Where("item_name ILIKE ?", "%"+s.Query+"%")
}
