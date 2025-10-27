package borrowq

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/configs"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories"
	borrowlog "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrow_log"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrowq"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
)

type BorrowQueueService interface {
	Enqueue(ctx context.Context, request requests.CreateBorrowQueueRequest) error
}

type borrowQueueService struct {
	cfg       *configs.Config
	bqRepo    borrowq.BorrowQueueRepository
	txManager repositories.TransactionManager
	borrowlog borrowlog.Repository
}

func NewBorrowQueueService(cfg *configs.Config, bqRepo borrowq.BorrowQueueRepository, txManager repositories.TransactionManager, borrowlog borrowlog.Repository) BorrowQueueService {
	return &borrowQueueService{
		cfg:       cfg,
		bqRepo:    bqRepo,
		txManager: txManager,
		borrowlog: borrowlog,
	}
}

// Entry implements BorrowQueueService.
func (b *borrowQueueService) Enqueue(ctx context.Context, request requests.CreateBorrowQueueRequest) error {
	userID, err := uuid.Parse(request.UserID)

	if err != nil {
		return err
	}

	return b.txManager.Do(ctx, func(ctx context.Context) error {
		b.borrowlog.CreateBorrowLog(ctx, models.BorrowLog{
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
			UserID: userID,
			ItemID: request.ItemID,
			CreatedAt: utils.BangkokNow(),
			BorrowedAt: nil,
		}

		err = b.bqRepo.Enqueue(ctx, data)
		if err != nil {
			return err
		}

		return nil
	})
}
