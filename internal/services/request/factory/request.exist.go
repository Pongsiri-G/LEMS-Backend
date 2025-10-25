package factory

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	ItemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	Minio "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/minio"
	RequestRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/request"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type existRequest struct {
	requestRepo RequestRepo.Repository
	itemRepo    ItemRepo.Repository
	minioRepo   Minio.Repository
	request     *models.Request
}

func NewExistRequestFactory(
	requestRepo RequestRepo.Repository,
	itemRepo ItemRepo.Repository,
	minioRepo Minio.Repository,
	request *models.Request,
) Requestable {
	return &existRequest{
		requestRepo: requestRepo,
		itemRepo:    itemRepo,
		minioRepo:   minioRepo,
		request:     request,
	}
}

// CreateRequest implements Requestable.
func (e *existRequest) CreateRequest(ctx context.Context, req requests.CreateRequest) error {
	itemID, err := uuid.Parse(*req.ItemID)
	if err != nil {
		return exceptions.ErrInvalidUUID
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return exceptions.ErrInvalidUUID
	}
	item, err := e.itemRepo.GetItemByID(ctx, itemID)
	if err != nil {
		log.Error().Err(err).Msg("failed to find item by ID")
		return err
	}

	if item == nil {
		return exceptions.ErrItemNotFound
	}

	request := models.Request{
		RequestID:          uuid.New(),
		UserID:             userID,
		RequestType:        req.RequestType,
		ItemID:             itemID,
		RequestImageURL:    new(string),
		RequestDescription: req.RequestDescription,
		CreatedAt:          utils.BangkokNow(),
		UpdatedAt:          utils.BangkokNow(),
	}

	return e.requestRepo.CreateRequest(ctx, &request)
}

// EditRequest implements Requestable.
func (e *existRequest) EditRequest(ctx context.Context, req requests.EditRequest) error {
	if req.RequestDescription != nil {
		e.request.RequestDescription = *req.RequestDescription
	}

	if req.ImageURL != nil {
		if e.request.RequestImageURL != nil {
			bucket, obj, err := utils.ExtractUrl(*e.request.RequestImageURL)
			if err != nil {
				log.Error().Err(err).Msg("failed to extract URL")
				return err
			}

			if err := e.minioRepo.DeleteImage(ctx, bucket, obj); err != nil {
				log.Error().Err(err).Msg("failed to delete old image from minio")
				return err
			}
		}

		e.request.RequestImageURL = req.ImageURL
	}

	err := e.requestRepo.EditRequest(ctx, e.request)
	if err != nil {
		log.Error().Err(err).Msg("failed to edit request")
		return err
	}

	return nil
}
