package responses

import (
	"time"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
)

type UserResponse struct {
	UserID         string             `json:"userId"`
	UserFullName   string             `json:"userFullName"`
	UserEmail      string             `json:"userEmail"`
	UserPhone      string             `json:"userPhone"`
	UserRole       enums.UserRole     `json:"userRole"`
	UserStatus     enums.UserStatus   `json:"userStatus"`
	UserProfileURL string             `json:"userProfileUrl"`
	AuthProvider   enums.AuthProvider `json:"authProvider"`
	LastLoggedIn   *time.Time         `json:"lastLoggedIn"`
}
