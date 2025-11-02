package borrowq

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories"
	borrowlog "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrow_log"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrowq"
	ItemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type BorrowQueueService interface {
	GetFrontQueue(ctx context.Context, itemID uuid.UUID) (*responses.BorrowQueueResponse, error)
	Enqueue(ctx context.Context, request requests.CreateBorrowQueueRequest) error
}

type borrowQueueService struct {
	cfg       *configs.Config
	bqRepo    borrowq.BorrowQueueRepository
	txManager repositories.TransactionManager
	borrowlog borrowlog.Repository
	itemRepo  ItemRepo.Repository
}

func NewBorrowQueueService(cfg *configs.Config, bqRepo borrowq.BorrowQueueRepository, txManager repositories.TransactionManager, borrowlog borrowlog.Repository, itemRepo ItemRepo.Repository) BorrowQueueService {
	return &borrowQueueService{
		cfg:       cfg,
		bqRepo:    bqRepo,
		txManager: txManager,
		borrowlog: borrowlog,
		itemRepo:  itemRepo,
	}
}

// Entry implements BorrowQueueService.
func (b *borrowQueueService) Enqueue(ctx context.Context, request requests.CreateBorrowQueueRequest) error {
	userID, err := uuid.Parse(request.UserID)

	if err != nil {
		return err
	}

	// only 1 enqueue per item
	existing, err := b.bqRepo.GetMemberByUserAndItem(ctx, request.ItemID.String(), request.UserID)
	if err != nil{
		return err
	}

	if existing != nil {
		return exceptions.ErrUserAlreadyBorrowq
	}

	return b.txManager.Do(ctx, func(ctx context.Context) error {
		b.borrowlog.CreateBorrowLogTx(ctx, models.BorrowLog{
			BorrowID:     uuid.New(),
			UserID:       userID,
			ItemID:       request.ItemID,
			BorrowStatus: enums.StatusWaiting,
			BorrowDate:   utils.BangkokNow(),
			ReturnDate:   nil,
			CreatedAt:    utils.BangkokNow(),
			UpdatedAt:    utils.BangkokNow(),
		})

		data := &models.BorrowQueue{
			UserID:     userID,
			ItemID:     request.ItemID,
			CreatedAt:  utils.BangkokNow(),
			BorrowedAt: nil,
		}

		err = b.bqRepo.Enqueue(ctx, data)
		if err != nil {
			return err
		}

		return nil
	})
}

// GetFrontQueue implements BorrowQueueService.
func (b *borrowQueueService) GetFrontQueue(ctx context.Context, itemID uuid.UUID) (*responses.BorrowQueueResponse, error) {
	item, err := b.itemRepo.GetItemByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, exceptions.ErrItemNotFound
	}

	queue, err := b.bqRepo.PeekOldest(ctx, itemID.String())
	if err != nil {
		log.Error().Err(err).Msg("failed to peek oldest borrow queue")
		return nil, err
	}

	if queue == nil {
		return nil, nil
	}

	resp := &responses.BorrowQueueResponse{
		QueueID:    queue.QueueID.String(),
		UserID:     queue.UserID.String(),
		ItemID:     queue.ItemID.String(),
		CreatedAt:  queue.CreatedAt,
		BorrowedAt: queue.BorrowedAt,
	}

	return resp, nil
}
