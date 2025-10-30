package responses

import (
	"time"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/google/uuid"
)

type ItemResponseBorrow struct {
	ID              uuid.UUID        `json:"id"`
	Name            string           `json:"name"`
	Description     *string          `json:"desc"`
	PictureURL      *string          `json:"picture_url"`
	Status          enums.ItemStatus `json:"status"`
	Quantity        int              `json:"quantity"`
	CurrentQuantity int              `json:"current_quantity"`
	Prerequisites   *[]ItemResponse  `json:"prerequisite,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	BorrowID        uuid.UUID        `json:"borrow_id"`
}
