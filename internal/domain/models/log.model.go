package models

import (
	"time"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/google/uuid"
)

type Log struct {
	LogID      uuid.UUID     `db:"log_id" gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID     `db:"user_id" gorm:"type:uuid;foreignKey:UserID;references:UserID"`
	LogType    enums.LogType `db:"log_type"`
	LogMessage *string       `db:"log_message"`
	CreatedAt  time.Time     `db:"created_at"`
}
