package models

import (
	"time"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/google/uuid"
)

type User struct {
	UserID         uuid.UUID          `db:"user_id" json:"user_id"`
	UserFullName   string             `db:"user_fullname" json:"user_fullname"`
	UserEmail      string             `db:"user_email" json:"user_email"`
	UserPhone      string             `db:"user_phone" json:"user_phone"`
	UserPassword   string             `db:"user_password" json:"-"` // hide in JSON
	UserRole       enums.UsesrRole    `db:"user_role" json:"user_role"`
	UserStatus     enums.UserStatus   `db:"user_status" json:"user_status"`
	UserProfileURL string             `db:"user_profile_url" json:"user_profile_url"`
	AuthProvider   enums.AuthProvider `db:"auth_provider" json:"auth_provider"`
	CreatedAt      time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `db:"updated_at" json:"updated_at"`
	LastLoggedIn   *time.Time         `db:"last_logged_in" json:"last_logged_in,omitempty"`
}
