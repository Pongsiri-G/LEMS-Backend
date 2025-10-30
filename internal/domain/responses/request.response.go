package responses

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/google/uuid"
)

type GetAllRequestsResponse struct {
	RequestID          uuid.UUID              `json:"request_id"`
	RequestItemName    string                 `json:"request_item_name"`
	RequestDescription string                 `json:"request_description"`
	ItemID             uuid.UUID              `json:"item_id"`
	ItemRequest        *ItemRequestedResponse `json:"item_requested"`
	RequestType        enums.RequestType      `json:"request_type"`
	RequestStatus      enums.RequestStatus    `json:"request_status"`
	RequestImageURL    *string                `json:"request_image_url"`
	RequestCreatedBy   string                 `json:"created_by"`
	RequestCreatedDate string                 `json:"created_date"`
	RequestUpdatedDate string                 `json:"updated_date"`
}

type ItemRequestedResponse struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}
