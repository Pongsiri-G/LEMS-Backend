package factory

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	ItemRequested "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item_requested"
	Minio "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/minio"
	RequestRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/request"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type withdrawRequestFactory struct {
	requestRepo       RequestRepo.Repository
	itemRequestedRepo ItemRequested.Repository
	minioRepo         Minio.Repository
	request           *models.Request
	userID            *uuid.UUID
}

func NewWithdrawRequestFactory(
	requestRepo RequestRepo.Repository,
	itemRequestedRepo ItemRequested.Repository,
	minioRepo Minio.Repository,
	request *models.Request,
	userID *uuid.UUID,
) Requestable {
	return &withdrawRequestFactory{
		requestRepo:       requestRepo,
		itemRequestedRepo: itemRequestedRepo,
		minioRepo:         minioRepo,
		request:           request,
		userID:            userID,
	}
}

// CreateRequest implements Requestable.
func (w *withdrawRequestFactory) CreateRequest(ctx context.Context, req requests.CreateRequest) error {
	if req.ItemID != nil {
		log.Error().Msg("request type does not expect item ID")
		return exceptions.ErrRequestNotExpectItemID
	}
	itemRequested := models.ItemRequested{
		ID:          uuid.New(),
		Name:        req.Item.Name,
		Description: req.Item.Description,
		Type:        req.Item.Type,
		UserID:      *w.userID,
		Quantity:    req.Item.Quantity,
		Price:       req.Item.Price,
	}

	if err := w.itemRequestedRepo.CreateItemRequested(ctx, &itemRequested); err != nil {
		log.Error().Err(err).Msg("failed to create item requested")
		return err
	}

	request := models.Request{
		RequestID:          uuid.New(),
		UserID:             *w.userID,
		RequestType:        req.RequestType,
		RequestDescription: req.RequestDescription,
		ItemID:             itemRequested.ID,
		RequestStatus:      enums.RequestStatusPending,
		RequestImageURL:    req.ImageURL,
		CreatedAt:          utils.BangkokNow(),
		UpdatedAt:          utils.BangkokNow(),
	}

	if err := w.requestRepo.CreateRequest(ctx, &request); err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return err
	}
	return nil
}

// EditRequest implements Requestable.
func (w *withdrawRequestFactory) EditRequest(ctx context.Context, req requests.EditRequest) error {

	item, err := w.itemRequestedRepo.FindByID(ctx, w.request.ItemID)
	if err != nil {
		log.Error().Err(err).Msg("failed to find item requested by ID")
		return err
	}

	if item == nil {
		return exceptions.ErrItemNotFound
	}

	if req.RequestDescription != nil {
		w.request.RequestDescription = *req.RequestDescription
	}

	if req.ImageURL != nil {
		if w.request.RequestImageURL != nil {
			// delete old image from minio
			bucketName, objectName, err := utils.ExtractUrl(*w.request.RequestImageURL)
			if err != nil {
				log.Error().Err(err).Msg("failed to extract URL")
				return err
			}
			if err := w.minioRepo.DeleteImage(ctx, bucketName, objectName); err != nil {
				log.Error().Err(err).Msg("failed to delete old image from minio")
				return err
			}
		}
		w.request.RequestImageURL = req.ImageURL
	}

	if req.ItemPrice != nil {
		item.Price = *req.ItemPrice
	}

	if req.ItemQuantity != nil {
		item.Quantity = *req.ItemQuantity
	}

	if req.ItemType != nil {
		item.Type = *req.ItemType
	}

	if req.ItemDescription != nil {
		item.Description = *req.ItemDescription
	}

	if req.ItemName != nil {
		item.Name = *req.ItemName
	}

	w.request.UpdatedAt = utils.BangkokNow()
	if err := w.requestRepo.EditRequest(ctx, w.request); err != nil {
		log.Error().Err(err).Msg("failed to edit request")
		return err
	}

	if err := w.itemRequestedRepo.EditItemRequested(ctx, item); err != nil {
		log.Error().Err(err).Msg("failed to edit item requested")
		return err
	}

	return nil

}
