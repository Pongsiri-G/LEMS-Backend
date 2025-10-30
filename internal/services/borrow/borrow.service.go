package borrow

import (
	"context"
	"fmt"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/events"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	borrowRepository "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrow_log"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrowq"
	itemRepository "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	itemsetRepository "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item_set"
	logsystem "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/log"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/item/factory"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/noti"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Return(ctx context.Context, req *requests.ReturnRequest) error
	Borrow(ctx context.Context, req *requests.BorrowRequest) error
}

type service struct {
	borrowRepo  borrowRepository.Repository
	itemRepo    itemRepository.Repository
	itemSetRepo itemsetRepository.Repository
	logRepo     logsystem.Repository
	events      noti.Subject
	bqRepo      borrowq.BorrowQueueRepository
}

func NewBorrowService(
	borrowRepo borrowRepository.Repository,
	itemRepo itemRepository.Repository,
	itemSetRepo itemsetRepository.Repository,
	logRepo logsystem.Repository,
	events noti.Subject,
	bqRepo borrowq.BorrowQueueRepository,
) Service {
	return &service{
		borrowRepo:  borrowRepo,
		itemRepo:    itemRepo,
		itemSetRepo: itemSetRepo,
		logRepo:     logRepo,
		events:      events,
		bqRepo:      bqRepo,
	}
}

func (s *service) Borrow(ctx context.Context, req *requests.BorrowRequest) error {
	var borrowFactory factory.Borrowable
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return exceptions.ErrInvalidUUID
	}
	itemID, err := uuid.Parse(req.ItemID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return exceptions.ErrInvalidUUID
	}

	item, err := s.itemRepo.GetItemByID(ctx, itemID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get item by id")
		return err
	}

	if item == nil {
		log.Error().Err(err).Msg("item not found")
		return exceptions.ErrItemNotFound
	}

	children, err := s.itemRepo.GetChildItemByParentID(ctx, item.ItemID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get child items")
		return err
	}
	if len(children) > 0 {
		borrowFactory = factory.NewChildItemBorrowable(s.itemRepo, s.borrowRepo, s.itemSetRepo, s.logRepo)
	} else {
		borrowFactory = factory.NewNormalItemBorrowable(s.itemRepo, s.borrowRepo, s.logRepo)
	}

	return borrowFactory.BorrowItem(ctx, userID, item, &children)

}

// Return implements Service.
func (s *service) Return(ctx context.Context, req *requests.ReturnRequest) error {
	var itemBorrowableFactory factory.Borrowable
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return exceptions.ErrInvalidUUID
	}
	borrowID, err := uuid.Parse(req.BorrowID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return exceptions.ErrInvalidUUID
	}

	borrow, err := s.borrowRepo.FindBorrowLogByUserIDAndBorrowID(ctx, userID, borrowID)
	if err != nil {
		log.Error().Err(err).Msg("failed to find borrow log")
		return err
	}

	if borrow == nil {
		log.Error().Err(err).Msg("borrow log not found")
		return exceptions.ErrBorrowLogNotFound
	}

	if borrow.BorrowParentID != nil {
		log.Error().Err(err).Msg("cannot return child item directly")
		return exceptions.ErrCannotReturnChildItemDirectly
	}

	children, err := s.borrowRepo.GetChildren(ctx, borrow.BorrowID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get child borrow logs")
		return err
	}

	if len(children) > 0 {
		itemBorrowableFactory = factory.NewChildItemBorrowable(s.itemRepo, s.borrowRepo, s.itemSetRepo, s.logRepo)
	} else {
		itemBorrowableFactory = factory.NewNormalItemBorrowable(s.itemRepo, s.borrowRepo, s.logRepo)
	}

	err = s.noitification(ctx, borrow.ItemID)
	if err != nil {
		return err
	}

	return itemBorrowableFactory.ReturnItem(ctx, borrow, &children)
}

func (s *service) noitification(ctx context.Context, itemID uuid.UUID) error {
	item, err := s.itemRepo.GetItemByID(ctx, itemID)
	if err != nil {
		return err
	}

	if item == nil {
		return exceptions.ErrItemNotFound
	}

	// Process queue (greedy: only head; call repeatedly or loop)
	next, err := s.bqRepo.PeekOldest(ctx, itemID.String())
	if err != nil {
		return err
	}

	if next == nil {
		return nil
	}

	// parsedUUID, err := uuid.Parse(itemID)
	// if err != nil {
	// 	return err
	// }

	// cur, _ := s.itemRepo.GetItemByID(ctx, parsedUUID)

	// mark the borrow request to READY (simple lookup not shown here; typically tie queue ID to borrow ID)
	// You can add a find-by(user,equipment,WAITING) here in a real repo
	err = s.bqRepo.Dequeue(ctx, next.QueueID)
	if err != nil {
		return err
	}

	// if s.events == nil {
	// 	continue
	// }

	s.events.Notify(events.Event{
		Type: events.ItemAvaliable,
		Payload: map[string]interface{}{
			"userId":      next.UserID.String(),
			"message":     fmt.Sprintf("Your requested equipment (%s) is ready for pickup", item.ItemName),
		},
	
	})	

	return nil
}
