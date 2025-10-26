package requests

import "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"

type CreateItemRequest struct {
	Name         string            `json:"name" validate:"required,max=100"`
	Description  *string           `json:"description" validate:"max=500"`
	ImageURL     *string           `json:"image_url"`
	Quantity     int               `json:"quantity" validate:"gte=0"`
	Status       *enums.ItemStatus `json:"status" validate:"omitempty"`
	Prerequisite *[]string         `json:"prerequisite" validate:"omitempty,max=500"`
	Tags         *[]string         `json:"tags" validate:"omitempty"`
}