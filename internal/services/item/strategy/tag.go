package strategy

import (
	"context"
	"strings"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
)

type TagStrategy struct {
	repo item.Repository
	data []string
}

func (f *TagStrategy) InitFilter(r item.Repository) {
	f.repo = r
}

func (f TagStrategy) Filter(ctx context.Context) ([]models.Item, error) {
	var tags []string
	unique := map[string]struct{}{}
	for _, tag := range f.data {
		for _, t := range strings.Split(tag, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				unique[t] = struct{}{}
			}
		}
	}

	for tag := range unique {
		tags = append(tags, tag)
	}

	f.data = tags

	return f.repo.GetByTags(ctx, f.data)
}
