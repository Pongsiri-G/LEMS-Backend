package strategy

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
)

type NameSearchStrategy struct {
	repo item.Repository
	data any
}

func (s *NameSearchStrategy) Init(r item.Repository) {
	s.repo = r
}

func (s NameSearchStrategy) Search(ctx context.Context) ([]models.Item, error) {
	return s.repo.GetByName(ctx, s.data.(string))
}