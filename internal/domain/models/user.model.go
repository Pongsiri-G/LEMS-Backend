package models

import (
	"time"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/google/uuid"
)

type User struct {
	UserID         uuid.UUID          `gorm:"type:uuid;primaryKey" json:"user_id"`
	UserFullName   string             `gorm:"type:varchar(50);not null" json:"user_fullname"`
	UserEmail      string             `gorm:"type:varchar(50);unique;not null" json:"user_email"`
	UserPhone      string             `gorm:"type:varchar(10)" json:"user_phone"`
	UserPassword   string             `gorm:"type:varchar(100);not null" json:"-"` // hidden in JSON
	UserRole       enums.UserRole     `gorm:"type:varchar(10);not null;default:'USER'" json:"user_role"`
	UserStatus     enums.UserStatus   `gorm:"type:varchar(20);not null;default:'PENDING'" json:"user_status"`
	UserProfileURL string             `gorm:"type:text" json:"user_profile_url"`
	AuthProvider   enums.AuthProvider `gorm:"type:varchar(50);not null;default:'LOCAL'" json:"auth_provider"`
	CreatedAt      time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      *time.Time         `gorm:"index" json:"deleted_at,omitempty"`
	LastLoggedIn   *time.Time         `json:"last_logged_in"`

	BorrowQueues []BorrowQueue `gorm:"foreignKey:UserID"`
}

func (User) TableName() string { return "users" }
