package item

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Repository interface {
	CreateItem(ctx context.Context, item *models.Item) error
	GetItemByID(ctx context.Context, itemID uuid.UUID) (*models.Item, error)
	UpdateItem(ctx context.Context, item *models.Item) error
	GetAll(ctx context.Context) ([]models.Item, error)
	GetMyBorrow(ctx context.Context, userID uuid.UUID) ([]models.ItemBorrow, error)
	GetChildItemByParentID(ctx context.Context, itemID uuid.UUID) ([]models.Item, error)
	GetAvailable(ctx context.Context) ([]models.Item, error)
	GetByTags(ctx context.Context, tags []string) ([]models.Item, error)
	GetByName(ctx context.Context, name string) ([]models.Item, error)
	SearchItems(ctx context.Context, strategies []SearchStrategy) ([]models.Item, error)
	DeleteItem(ctx context.Context, itemID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

type SearchStrategy interface {
	Apply(db *gorm.DB) *gorm.DB
}

type SearchStrategyMap struct {
	Tags   []string
	Name   string
	Status string
	User   string
}

// Constuctor
func NewItemRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetChildItemByParentID(ctx context.Context, itemID uuid.UUID) ([]models.Item, error) {
	var itemSets []models.ItemSets
	items := make([]models.Item, 0)
	err := r.db.Where("parent_item_id = ?", itemID).Find(&itemSets).Error
	if err != nil {
		log.Error().Err(err).Msg("can't get item sets from database")
		return []models.Item{}, err
	}
	if len(itemSets) <= 0 {
		return items, nil
	}
	for _, itemSet := range itemSets {
		var item models.Item
		if err := r.db.First(&item, itemSet.ChildItemID).Error; err != nil {
			log.Error().Err(err).Msg("can't get item from parent id that get from item sets")
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *repository) UpdateItem(ctx context.Context, item *models.Item) error {
	return r.db.Save(item).Error
}

func (r *repository) GetItemByID(ctx context.Context, itemID uuid.UUID) (*models.Item, error) {
	var item models.Item
	err := r.db.WithContext(ctx).Where("item_id = ?", itemID).First(&item).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Err(err).Msg("item not found")
			return nil, nil
		}
		log.Error().Err(err).Msg("failed to get item by id")
		return nil, err
	}

	return &item, nil

}

// CreateItem implements Repository.
func (r *repository) CreateItem(ctx context.Context, item *models.Item) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *repository) GetAll(ctx context.Context) ([]models.Item, error) {
	var items []models.Item
	err := r.db.Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) GetMyBorrow(ctx context.Context, userID uuid.UUID) ([]models.ItemBorrow, error) {
	var items []models.ItemBorrow
	err := r.db.
		Table("items").
		Select("items.item_id, items.item_name, borrow_logs.borrow_id, items.item_description, items.item_picture_url, items.item_status, items.item_quantity, items.item_created_at, items.item_updated_at, items.item_current_quantity").
		Joins("JOIN borrow_logs ON items.item_id::uuid = borrow_logs.item_id").
		Where("borrow_logs.user_id = ? AND borrow_logs.borrow_status = ?", userID, "BORROWED").
		Find(&items).Error

	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) GetAvailable(ctx context.Context) ([]models.Item, error) {
	var items []models.Item
	err := r.db.Where("item_status=?", "AVAILABLE").Find(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *repository) GetByTags(ctx context.Context, tags []string) ([]models.Item, error) {
	var items []models.Item

	err := r.db.Table("items AS i").
		Select("i.*").
		Joins("JOIN item_tags it ON i.item_id = it.item_id").
		Joins("JOIN tags t ON it.tag_id = t.tag_id").
		Where("t.tag_name IN (?)", tags).
		Scan(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *repository) GetByName(ctx context.Context, name string) ([]models.Item, error) {
	var items []models.Item
	err := r.db.Where("item_name ILIKE ?", "%"+name+"%").Find(&items).Error

	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) SearchItems(ctx context.Context, strategies []SearchStrategy) ([]models.Item, error) {
	var items []models.Item
	log.Info().Msgf("Searching items with %d strategies", len(strategies))
	db := r.db.Model(&models.Item{})

	for _, strategy := range strategies {
		db = strategy.Apply(db)
	}

	err := db.Find(&items).Error
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("Found %d items matching search criteria", len(items))
	return items, nil
}

// DeleteItem implements Repository.
func (r *repository) DeleteItem(ctx context.Context, itemID uuid.UUID) error {
	result := r.db.Where("item_id = ?", itemID).Delete(&models.Item{})
	if result.RowsAffected == 0 {
		return nil
	}
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("failed to delete item")
		return result.Error
	}

	return nil
}
