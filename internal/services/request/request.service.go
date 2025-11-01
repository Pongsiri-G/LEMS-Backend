package request

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	ItemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	ItemRequested "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item_requested"
	Minio "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/minio"
	RequestRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/request"
	User "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/request/factory"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	GetRequests(ctx context.Context, userID *uuid.UUID) ([]responses.GetAllRequestsResponse, error)

	CreateRequest(ctx context.Context, userID uuid.UUID, req requests.CreateRequest) error
	EditRequest(ctx context.Context, req requests.EditRequest) error
	ExportRequests(ctx context.Context, req requests.ExportRequests) error
	ChangeRequestStatus(ctx context.Context, requestID string, status enums.RequestStatus) error
}

type service struct {
	requestRepo       RequestRepo.Repository
	itemRequestedRepo ItemRequested.Repository
	itemRepo          ItemRepo.Repository
	userRepo          User.Repository
	minioRepo         Minio.Repository
}

func NewRequestService(
	requestRepo RequestRepo.Repository,
	itemRequestedRepo ItemRequested.Repository,
	itemRepo ItemRepo.Repository,
	minioRepo Minio.Repository,
	userRepo User.Repository,
) Service {
	return &service{
		requestRepo:       requestRepo,
		itemRequestedRepo: itemRequestedRepo,
		itemRepo:          itemRepo,
		minioRepo:         minioRepo,
		userRepo:          userRepo,
	}
}

// CreateRequest implements Service.
func (s *service) CreateRequest(ctx context.Context, userID uuid.UUID, req requests.CreateRequest) error {

	var requestFactory factory.Requestable
	ok := enums.IsValidRequestType(req.RequestType)
	if !ok {
		log.Error().Msg("invalid request type")
		return exceptions.ErrRequestInvalidRequestType
	}
	if req.RequestType == enums.RequestTypeRequest {
		if req.Item == nil {
			log.Error().Msg("item requested is nil for request type 'request'")
			return exceptions.ErrRequestItemInvalid
		}
		requestFactory = factory.NewWithdrawRequestFactory(s.requestRepo, s.itemRequestedRepo, s.minioRepo, nil, userID)
	} else {
		if req.ItemID == nil {
			return exceptions.ErrRequestItemIDInvalid
		}
		requestFactory = factory.NewExistRequestFactory(s.requestRepo, s.itemRepo, s.minioRepo, userID, nil)
	}

	return requestFactory.CreateRequest(ctx, req)
}

func (s *service) EditRequest(ctx context.Context, req requests.EditRequest) error {
	var requestFactory factory.Requestable

	requestID, err := uuid.Parse(req.RequestID)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse request ID")
		return exceptions.ErrInvalidUUID
	}
	request, err := s.requestRepo.FindByID(ctx, requestID)
	if err != nil {
		log.Error().Err(err).Msg("failed to find request by ID")
		return err
	}

	if request == nil {
		return exceptions.ErrRequestNotFound
	}

	ok := enums.IsValidRequestType(request.RequestType)
	if !ok {
		log.Error().Msg("invalid request type")
		return exceptions.ErrRequestInvalidRequestType
	}

	if request.RequestType == enums.RequestTypeRequest {
		requestFactory = factory.NewWithdrawRequestFactory(s.requestRepo, s.itemRequestedRepo, s.minioRepo, request, uuid.UUID{})
	} else {
		requestFactory = factory.NewExistRequestFactory(s.requestRepo, s.itemRepo, s.minioRepo, uuid.UUID{}, request)
	}

	return requestFactory.EditRequest(ctx, req)
}

// GetRequests implements Service.
func (s *service) GetRequests(ctx context.Context, userID *uuid.UUID) ([]responses.GetAllRequestsResponse, error) {
	var requestsData []models.Request
	var err error
	if userID != nil {
		requestsData, err = s.requestRepo.GetRequestsByUserID(ctx, *userID, nil, nil)
	} else {
		requestsData, err = s.requestRepo.GetRequests(ctx, nil, nil)
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to get requests")
		return nil, err
	}

	var response []responses.GetAllRequestsResponse = make([]responses.GetAllRequestsResponse, 0)
	for _, req := range requestsData {
		user, err := s.userRepo.FindByID(ctx, req.UserID.String())
		if err != nil {
			log.Error().Err(err).Msg("failed to find user by ID")
			return nil, err
		}

		if user == nil {
			log.Error().Msg("user not found for user ID: " + req.UserID.String())
			return nil, exceptions.ErrUserNotFound
		}

		var res = responses.GetAllRequestsResponse{
			RequestID:          req.RequestID,
			RequestType:        req.RequestType,
			RequestStatus:      req.RequestStatus,
			RequestImageURL:    req.RequestImageURL,
			RequestCreatedBy:   user.UserFullName,
			Quantity:           req.Quantity,
			RequestCreatedDate: utils.ToStringDateTime(req.CreatedAt),
			RequestUpdatedDate: utils.ToStringDateTime(req.UpdatedAt),
		}

		if req.RequestType == enums.RequestTypeRequest {
			itemRequested, err := s.itemRequestedRepo.FindByID(ctx, req.ItemID)
			if err != nil {
				log.Error().Err(err).Msg("failed to get item requested by ID")
				return nil, err
			}
			if itemRequested == nil {
				log.Error().Msg("item requested not found for item ID: " + req.ItemID.String())
				return nil, exceptions.ErrRequestedItemNotFound
			}
			res.RequestItemName = itemRequested.Name
			res.ItemID = itemRequested.ID
			res.RequestDescription = req.RequestDescription
			res.ItemRequest = &responses.ItemRequestedResponse{
				Name:        itemRequested.Name,
				Description: itemRequested.Description,
				Type:        itemRequested.Type,
				Quantity:    itemRequested.Quantity,
				Price:       itemRequested.Price,
			}
		} else {
			item, err := s.itemRepo.GetItemByID(ctx, req.ItemID)
			if err != nil {
				log.Error().Err(err).Msg("failed to get item by ID")
				return nil, err
			}
			if item == nil {
				log.Error().Msg("item not found for item ID: " + req.ItemID.String())
				return nil, exceptions.ErrItemNotFound
			}

			res.RequestItemName = item.ItemName
			res.RequestDescription = req.RequestDescription
			res.ItemID = item.ItemID
		}

		response = append(response, res)

	}

	return response, nil
}

// ChangeRequestStatus implements Service.
func (s *service) ChangeRequestStatus(ctx context.Context, requestID string, status enums.RequestStatus) error {
	requestUUID, err := uuid.Parse(requestID)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse request ID")
		return exceptions.ErrInvalidUUID
	}

	request, err := s.requestRepo.FindByID(ctx, requestUUID)
	if err != nil {
		log.Error().Err(err).Msg("failed to find request by ID")
		return err
	}

	request.UpdatedAt = utils.BangkokNow()
	request.RequestStatus = status
	return s.requestRepo.EditRequest(ctx, request)
}

// ExportRequests implements Service.
func (s *service) ExportRequests(ctx context.Context, req requests.ExportRequests) error {
	panic("unimplement")
}
