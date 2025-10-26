package repository

import "gorm.io/gorm"

type UserSearch struct {
	Query string
}

func (s UserSearch) Apply(db *gorm.DB) *gorm.DB {
	if s.Query == "" {
		return db
	}
	return db.Joins("JOIN borrow_logs bl ON items.item_id = bl.item_id").
			Where("bl.user_id = ? AND bl.borrow_status = 'BORROWED'", s.Query)
}