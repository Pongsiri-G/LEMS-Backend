package itemset

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository interface {
	CreateItemSet(parentID, childID uuid.UUID) error
	GetChildItemByParentID(parentID uuid.UUID) ([]models.ItemSets, error)
	DeleteItemSet(parentID, childID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewItemSetRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// DeleteItemSet implements Repository.
func (r *repository) DeleteItemSet(parentID uuid.UUID, childID uuid.UUID) error {
	result := r.db.Where("parent_item_id = ? AND child_item_id = ?", parentID, childID).Delete(&models.ItemSets{})
	if result.RowsAffected == 0 {
		return nil
	}
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("failed to delete item set")
		return result.Error
	}
	return nil
}

// CreateItemSet implements Repository.
func (r *repository) CreateItemSet(parentID, childID uuid.UUID) error {
	return r.db.Create(&models.ItemSets{
		ParentItemID: parentID,
		ChildItemID:  childID,
	}).Error
}

// GetChildItemByParentID implements Repository.
func (r *repository) GetChildItemByParentID(parentID uuid.UUID) ([]models.ItemSets, error) {
	var itemSets []models.ItemSets
	if err := r.db.Where("parent_item_id = ?", parentID).Find(&itemSets).Error; err != nil {
		log.Error().Err(err).Msg("failed to get child items by parent ID")
		return nil, err
	}
	return itemSets, nil
}
