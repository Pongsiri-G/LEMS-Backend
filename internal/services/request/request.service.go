package request

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	ItemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	ItemRequested "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item_requested"
	Minio "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/minio"
	RequestRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/request"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/request/factory"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	CreateRequestRequest(ctx context.Context, req requests.CreateRequest) error
}

type service struct {
	requestRepo       RequestRepo.Repository
	itemRequestedRepo ItemRequested.Repository
	itemRepo          ItemRepo.Repository
	minioRepo         Minio.Repository
}

func NewRequestService(
	requestRepo RequestRepo.Repository,
	itemRequestedRepo ItemRequested.Repository,
	itemRepo ItemRepo.Repository,
	minioRepo Minio.Repository,
) Service {
	return &service{
		requestRepo:       requestRepo,
		itemRequestedRepo: itemRequestedRepo,
		itemRepo:          itemRepo,
	}
}

// CreateRequestRequest implements Service.
func (s *service) CreateRequestRequest(ctx context.Context, req requests.CreateRequest) error {
	var requestFactory factory.Requestable
	if req.RequestType == enums.RequestTypeRequest {
		requestFactory = factory.NewWithdrawRequestFactory(s.requestRepo, s.itemRequestedRepo, s.minioRepo, nil)
	} else {
		requestFactory = factory.NewExistRequestFactory(s.requestRepo, s.itemRepo, s.minioRepo, nil)
	}

	return requestFactory.CreateRequest(ctx, req)
}

func (s *service) EditRequestRequest(ctx context.Context, req requests.EditRequest) error {
	var requestFactory factory.Requestable

	requestID, err := uuid.Parse(req.RequestID)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse request ID")
		return err
	}
	request, err := s.requestRepo.FindByID(ctx, requestID)
	if err != nil {
		log.Error().Err(err).Msg("failed to find request by ID")
		return err
	}

	if request.RequestType == enums.RequestTypeRequest {
		requestFactory = factory.NewWithdrawRequestFactory(s.requestRepo, s.itemRequestedRepo, s.minioRepo, request)
	} else {
		requestFactory = factory.NewExistRequestFactory(s.requestRepo, s.itemRepo, s.minioRepo, request)
	}

	return requestFactory.EditRequest(ctx, req)
}