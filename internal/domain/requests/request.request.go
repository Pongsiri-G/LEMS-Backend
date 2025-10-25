package requests

import "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"

type CreateRequest struct {
	UserID             string                `json:"user_id" validate:"required"`
	RequestType        enums.RequestType     `json:"request_type" validate:"required"`
	RequestDescription string                `json:"request_description" validate:"required"`
	ImageURL           string                `json:"image_url" validate:"required"`
	Item               *ItemRequestedRequest `json:"item_requested" validate:"omitempty,dive"`
	ItemID             *string               `json:"item_id" validate:"omitempty"`
}

type ItemRequestedRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Type        string  `json:"type" validate:"required"`
	UserID      string  `json:"user_id" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}

type EditRequest struct {
	RequestID          string  `json:"request_id" validate:"required"`
	RequestDescription *string `json:"request_description" validate:"required"`
	ImageURL           *string `json:"image_url" validate:"required"`
	// for item
	ItemName        *string  `json:"item_name" validate:"omitempty"`
	ItemDescription *string  `json:"item_description" validate:"omitempty"`
	ItemType        *string  `json:"item_type" validate:"omitempty"`
	ItemQuantity    *int     `json:"item_quantity" validate:"omitempty,min=1"`
	ItemPrice       *float64 `json:"item_price" validate:"omitempty,gt=0"`
}
