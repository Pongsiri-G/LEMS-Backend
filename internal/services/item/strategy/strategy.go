package strategy

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
)

type FilterStrategy interface {
	InitFilter(repo item.Repository)
	Filter(ctx context.Context) ([]models.Item, error)
}

func NewFilterMap(data []string) map[string]FilterStrategy {
	return map[string]FilterStrategy{
		"available": &AvailabilityFilter{data: data},
		"tags":      &TagStrategy{data: data},
	}
}
