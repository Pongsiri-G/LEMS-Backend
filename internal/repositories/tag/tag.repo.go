package tag

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository interface {
	GetTagsByItemID(ctx context.Context, itemID uuid.UUID) ([]models.Tag, error)
	GetAllTags(ctx context.Context) ([]models.Tag, error)

	CreateTag(ctx context.Context, tag *models.Tag) error
	AssignTagToItem(ctx context.Context, itemID uuid.UUID, tagID uuid.UUID) error
	FindAssignedTagsByItemIDAndTagID(ctx context.Context, itemID uuid.UUID, tagID uuid.UUID) (*models.ItemTag, error)
	GetTagByName(ctx context.Context, tagName string) (*models.Tag, error)
	GetTagByID(ctx context.Context, tagID uuid.UUID) (*models.Tag, error)
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

// GetAllTags implements Repository.
func (r *repository) GetAllTags(ctx context.Context) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Find(&tags).Error
	if err != nil {
		log.Error().Err(err).Msg("can't get tags from database")
		return []models.Tag{}, err
	}
	return tags, nil
}

// CreateTag implements Repository.
func (r *repository) CreateTag(ctx context.Context, tag *models.Tag) error {
	err := r.db.Create(tag).Error
	if err != nil {
		log.Error().Err(err).Msg("failed to create tag in database")
		return err
	}
	return nil
}

// AssignTagToItem implements Repository.
func (r *repository) AssignTagToItem(ctx context.Context, itemID uuid.UUID, tagID uuid.UUID) error {
	tag, err := r.FindAssignedTagsByItemIDAndTagID(ctx, itemID, tagID)
	if err != nil {
		return err
	}
	if tag != nil {
		return exceptions.ErrTagAlreadyAssigned
	}
	itemTag := models.ItemTag{
		ItemID: itemID,
		TagID:  tagID,
	}
	err = r.db.Create(&itemTag).Error
	if err != nil {
		log.Error().Err(err).Msg("failed to assign tag to item in database")
		return err
	}

	return nil
}

// GetTagByID implements Repository.
func (r *repository) GetTagByID(ctx context.Context, tagID uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.First(&tag, "tag_id = ?", tagID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error().Err(err).Msg("failed to get tag by ID from database")
		return nil, err
	}
	return &tag, nil
}

// GetTagByName implements Repository.
func (r *repository) GetTagByName(ctx context.Context, tagName string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.First(&tag, "tag_name = ?", tagName).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error().Err(err).Msg("failed to get tag by name from database")
		return nil, err
	}
	return &tag, nil
}

// FindAssignedTagsByItemIDAndTagID implements Repository.
func (r *repository) FindAssignedTagsByItemIDAndTagID(ctx context.Context, itemID uuid.UUID, tagID uuid.UUID) (*models.ItemTag, error) {
	var itemTag models.ItemTag
	err := r.db.First(&itemTag, "item_id = ? AND tag_id = ?", itemID, tagID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error().Err(err).Msg("failed to find assigned tag by item ID and tag ID from database")
		return nil, err
	}
	return &itemTag, nil
}
