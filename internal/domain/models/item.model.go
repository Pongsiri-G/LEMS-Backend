package models

import (
	"time"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	"github.com/google/uuid"
)

type ItemInterface interface {
	ToResponse() responses.ItemResponse
	FromRequest(req *requests.CreateItemRequest) error
}

type Item struct {
	ItemID              uuid.UUID        `db:"item_id" gorm:"primaryKey;type:uuid"`
	ItemName            string           `db:"item_name"`
	ItemDescription     *string          `db:"item_description"`
	ItemPictureURL      *string          `db:"item_picture_url"`
	ItemStatus          enums.ItemStatus `db:"item_status"`
	ItemQuantity        int              `db:"item_quantity"`
	ItemCurrentQuantity int              `db:"item_current_quantity"`
	ItemCreatedAt       time.Time        `db:"item_created_at"`
	ItemUpdatedAt       time.Time        `db:"item_updated_at"`
}

type ItemParentChild struct {
	ParentID uuid.UUID `db:"parent_id" gorm:"type:uuid"`
	ChildID  uuid.UUID `db:"child_id" gorm:"type:uuid"`
}

type ItemWithChildren struct {
	Item
	Children []Item
}

// FromRequest implements ItemInterface.
func (i *ItemWithChildren) FromRequest(req *requests.CreateItemRequest) error {
	i.ItemID = uuid.New()
	i.ItemName = req.Name
	i.ItemDescription = req.Description
	i.ItemPictureURL = req.ImageURL
	i.ItemQuantity = req.Quantity
	i.ItemCurrentQuantity = req.Quantity
	i.ItemCreatedAt = time.Now()
	i.ItemUpdatedAt = time.Now()

	if req.Status != nil {
		i.ItemStatus = *req.Status
	} else {
		i.ItemStatus = enums.ItemStatusAvailable
	}
	return nil
}

// ToResponse implements ItemInterface.
func (i *ItemWithChildren) ToResponse() responses.ItemResponse {
	panic("unimplemented")
}

func (Item) TableName() string { return "items" }
