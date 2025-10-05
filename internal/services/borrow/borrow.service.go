package borrow

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	borrowRepository "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrow_log"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Return(ctx context.Context, req *requests.ReturnRequest) error
}

type service struct {
	borrowRepo borrowRepository.Repository
}

func NewBorrowService(borrowRepo borrowRepository.Repository) Service {
	return &service{borrowRepo: borrowRepo}
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
