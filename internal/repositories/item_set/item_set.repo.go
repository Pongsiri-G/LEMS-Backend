package itemset

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateItemSet(parentID, childID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

// CreateItemSet implements Repository.
func (r *repository) CreateItemSet(parentID, childID uuid.UUID) error {
	return r.db.Create(&models.ItemSets{
		ParentItemID: parentID,
		ChildItemID:  childID,
	}).Error
}

func NewItemSetRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
