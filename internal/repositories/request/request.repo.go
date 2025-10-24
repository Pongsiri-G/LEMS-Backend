package request

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateRequest(request *models.Request) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateRequest implements Repository.
func (r *repository) CreateRequest(request *models.Request) error {
	return r.db.Create(request).Error
}
