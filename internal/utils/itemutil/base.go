package itemutil

import (
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
)

func ToResponses(items []models.Item) []responses.ItemResponse {
	var response []responses.ItemResponse
	for _, i := range items {
		r := ToResponse(i)
		response = append(response, r)
	}
	return response
}

func ToResponse(item models.Item) responses.ItemResponse {
	return responses.ItemResponse{
		ID:          item.ItemID,
		Name:        item.ItemName,
		Description: item.ItemDescription,
		PictureURL:  item.ItemPictureURL,
		Status:      item.ItemStatus,
		Quantity:    item.ItemQuantity,
		CreatedAt:   item.ItemCreatedAt,
		UpdatedAt:   item.ItemUpdatedAt,
	}
}
