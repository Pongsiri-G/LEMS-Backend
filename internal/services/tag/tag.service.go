package tag

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	tagRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/tag"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	GetTagsNameByItemID(ctx context.Context, itemID string) ([]responses.TagResponse, error)
	GetAllTags(ctx context.Context) ([]responses.TagResponse, error)
}

type tagService struct {
	repo tagRepo.Repository
}

func NewTagService(repo tagRepo.Repository) Service {
	return &tagService{repo: repo}
}

func (i *tagService) GetTagsNameByItemID(ctx context.Context, itemID string) ([]responses.TagResponse, error) {
	itemIDUUID, err := uuid.Parse(itemID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return []responses.TagResponse{}, ErrInvalidUUID
	}
	tags, err := i.repo.GetTagsByItemID(ctx, itemIDUUID)
	if err != nil {
		return []responses.TagResponse{}, err
	}
	response := make([]responses.TagResponse, 0)

	if len(tags) <= 0 {
		return response, err
	}

	for _, t := range tags {
		r := responses.TagResponse{
			TagID:    t.TagID,
			TagName:  t.TagName,
			TagColor: t.TagColor,
		}
		response = append(response, r)
	}
	return response, nil
}

// GetAllTags implements Service.
func (i *tagService) GetAllTags(ctx context.Context) ([]responses.TagResponse, error) {
	tags, err := i.repo.GetAllTags(ctx)
	if err != nil {
		return []responses.TagResponse{}, err
	}
	response := make([]responses.TagResponse, 0)

	if len(tags) <= 0 {
		return response, err
	}

	for _, t := range tags {
		r := responses.TagResponse{
			TagID:    t.TagID,
			TagName:  t.TagName,
			TagColor: t.TagColor,
		}
		response = append(response, r)
	}
	return response, nil
}
