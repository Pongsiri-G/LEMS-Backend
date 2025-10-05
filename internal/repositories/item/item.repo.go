package item

import (
	"context"
	"fmt"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"

)

type Repository interface {
	CreateItem(ctx context.Context, item *models.Items) error
	GetItemByID(ctx context.Context, itemID uuid.UUID) (*models.Items, error)
	UpdateItem(ctx context.Context, item *models.Items) error
	GetAll(ctx context.Context) ([]models.Items, error)
	GetMyBorrow(ctx context.Context, userID uuid.UUID) ([]models.Items, error)
}

type repository struct {
	db *gorm.DB
}

func (r *repository) UpdateItem(ctx context.Context, item *models.Items) error {
	return r.db.Save(item).Error
}

func (r *repository) GetItemByID(ctx context.Context, itemID uuid.UUID) (*models.Items, error) {
	var item models.Items
	err := r.db.WithContext(ctx).Where("item_id = ?", itemID).First(&item).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println(("ERR Record not found"))
			return nil, err
		}
		return nil, err
	}

	return &item, nil

}

func NewItemRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateItem implements Repository.
func (r *repository) CreateItem(ctx context.Context, item *models.Items) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *repository) GetAll(ctx context.Context) ([]models.Items, error) {
	var items []models.Items
	err := r.db.Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) GetMyBorrow(ctx context.Context, userID uuid.UUID) ([]models.Items, error) {
	var items []models.Items

	sub := r.db.Model(&models.BorrowLog{}).Select("item_id").Where("user_id::uuid = ? AND borrow_status = ?", userID, "BORROWED")

	err := r.db.Where("item_id IN (?)", sub).Find(&items).Error
	if err != nil {
		return nil, err
	}
	
	return items, nil
}