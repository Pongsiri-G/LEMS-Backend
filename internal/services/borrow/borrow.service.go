package borrow

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	borrowRepository "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrow_log"
	itemRepository "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Return(ctx context.Context, req *requests.ReturnRequest) error
	Borrow(ctx context.Context, req *requests.BorrowRequest) error
}

type service struct {
	borrowRepo borrowRepository.Repository
	itemRepo   itemRepository.Repository
}

func NewBorrowService(borrowRepo borrowRepository.Repository, itemRepo itemRepository.Repository) Service {
	return &service{borrowRepo: borrowRepo, itemRepo: itemRepo}
}

func (s *service) Borrow(ctx context.Context, req *requests.BorrowRequest) error {

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return ErrInvalidUUID
	}
	itemID, err := uuid.Parse(req.ItemID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return ErrInvalidUUID
	}
	borrowLog := models.BorrowLog{
		BorrowID:     uuid.New(),
		UserID:       userID,
		ItemID:       itemID,
		BorrowStatus: enums.StatusReturned,
		BorrowDate:   utils.BangkokNow(),
		ReturnDate:   nil,
		CreatedAt:    utils.BangkokNow(),
		UpdatedAt:    utils.BangkokNow(),
	}

	item, err := s.itemRepo.GetItemByID(ctx, itemID)
	if item.ItemQuantity-1 < 0 {
		log.Error().Err(err).Msg("quantity of the item less than zero")
		return ErrItemQuantityInSufficient
	}

	item.ItemQuantity -= 1
	updateErr := s.itemRepo.UpdateItem(ctx, item)
	if updateErr != nil {
		log.Error().Err(err).Msg("failed to update quentitiy of item")
		return ErrFailedToUpdateQuantity
	}

	return s.borrowRepo.CreateBorrowLog(ctx, borrowLog)
}

// Return implements Service.
func (s *service) Return(ctx context.Context, req *requests.ReturnRequest) error {
	userID, err := uuid.FromBytes([]byte(req.UserID))
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return ErrInvalidUUID
	}
	borrowID, err := uuid.FromBytes([]byte(req.StoreID))
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return ErrInvalidUUID
	}
	borrow, err := s.borrowRepo.FindBorrowLogByUserIDAndBorrowID(ctx, userID, borrowID)
	if err != nil {
		return err
	}

	borrow.BorrowStatus = enums.StatusReturned
	now := utils.BangkokNow()
	borrow.ReturnDate = &now
	borrow.UpdatedAt = now

	return s.borrowRepo.EditBorrowLog(ctx, borrow)
}
