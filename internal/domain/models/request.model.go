package models

import (
	"time"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/google/uuid"
)

type Request struct {
	RequestID          uuid.UUID         `db:"request_id" gorm:"type:uuid;primaryKey;"`
	UserID             uuid.UUID         `db:"user_id" gorm:"type:uuid;not null;"`
	RequestType        enums.RequestType `db:"request_type" gorm:"type:varchar(50);not null;"`
	RequestTitle       string            `db:"request_title" gorm:"not null;"`
	RequestDescription string            `db:"request_description" gorm:"type:text;not null;"`
	CreatedAt          time.Time         `db:"created_at" gorm:"type:timestamp;not null;"`
	UpdatedAt          time.Time         `db:"updated_at" gorm:"type:timestamp;not null;"`
}
