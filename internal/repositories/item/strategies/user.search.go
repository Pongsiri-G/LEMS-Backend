package repository

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserSearch struct {
	Query string
}

func (s UserSearch) Apply(db *gorm.DB) *gorm.DB {
	if s.Query == "" {
		return db
	}
	userID, err := uuid.Parse(s.Query)
	if err != nil {
		log.Error().Err(err).Msg("error some thing eiei")
		return db
	}
	return db.Joins("JOIN borrow_logs bl ON items.item_id::uuid = bl.item_id").
		Where("bl.user_id = ? AND bl.borrow_status = 'BORROWED'", userID)
}
