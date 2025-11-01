package requests

import "github.com/google/uuid"

type CreateBorrowQueueRequest struct {
	UserID string    
	ItemID uuid.UUID `json:"itemId" validate:"required"`
}
