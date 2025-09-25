package repositories

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"gorm.io/gorm"
)

type UserRepoistory interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepoistory {
	return &userRepository{
		db: db,
	}
}

// Create implements UserRepoistory.
func (u *userRepository) Create(ctx context.Context, user *models.User) error {
	panic("unimplemented")
}

// FindByEmail implements UserRepoistory.
func (u *userRepository) FindByEmail(ctx context.Context, email string) error {
	panic("unimplemented")
}
