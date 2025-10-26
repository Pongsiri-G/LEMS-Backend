package strategy

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
)

type AvailabilityFilter struct {
	repo item.Repository
	data []string
}

func (f *AvailabilityFilter) InitFilter(r item.Repository) {
	f.repo = r
}

func (f AvailabilityFilter) Filter(ctx context.Context) ([]models.Item, error) {
	return f.repo.GetAvailable(ctx)
}
