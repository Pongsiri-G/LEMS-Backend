package tag

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository interface {
	GetTagsByItemID(ctx context.Context, itemID uuid.UUID) ([]models.Tag, error)
}

type repository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetTagsByItemID(ctx context.Context, itemID uuid.UUID) ([]models.Tag, error) {
	var tagsItem []models.ItemTag
	tags := make([]models.Tag, 0)
	err := r.db.Where("item_id = ?", itemID).Find(&tagsItem).Error
	if err != nil {
		log.Error().Err(err).Msg("can't get item tags from database")
		return []models.Tag{}, err
	}
	if len(tagsItem) <= 0 {
		return tags, nil
	}
	for _, itemTag := range tagsItem {
		var tag models.Tag
		if err := r.db.First(&tag, itemTag.TagID).Error; err != nil {
			log.Error().Err(err).Msg("can't get item from item ID in item tag")
			continue
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
