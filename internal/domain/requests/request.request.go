package requests

type CreateRequestRequest struct {
}

type ItemRequestedRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Type        string  `json:"type" validate:"required"`
	UserID      string  `json:"user_id" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}

// type ItemRequested struct {
// 	ID          uuid.UUID `db:"item_requested_id" gorm:"primaryKey;type:uuid"`
// 	Name        string    `db:"item_requested_name" gorm:"not null"`
// 	Description string    `db:"item_requested_description" gorm:"not null"`
// 	Type        string    `db:"item_requested_type" gorm:"not null"`
// 	UserID      uuid.UUID `db:"user_id" gorm:"type:uuid;not null"`
// 	Quantity    int       `db:"quantity" gorm:"not null"`
// 	Price       float64   `db:"price" gorm:"not null"`
// }
