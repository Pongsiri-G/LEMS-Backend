package requests

//	type Item struct {
//		ItemID          uuid.UUID `db:"item_id" gorm:"primaryKey"`
//		ItemName        string    `db:"item_name"`
//		ItemDescription string    `db:"item_description"`
//		ItemPictureURL  string    `db:"item_picture_url"`
//		ItemStatus      string    `db:"item_status"`
//		ItemQuantity    int       `db:"item_quantity"`
//		ItemCreatedAt   time.Time `db:"item_created_at"`
//		ItemUpdatedAt   time.Time `db:"item_updated_at"`
//	}
type CreateItemRequest struct {
	Name        string   `json:"name" validate:"required,max=100"`
	Description *string  `json:"description" validate:"max=500"`
	Price       *float64 `json:"price" validate:"gte=0"`
	ImageURL    *string  `json:"image_url"`
	Quantity    int      `json:"quantity" validate:"gte=0"`
}
