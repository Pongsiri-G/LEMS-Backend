package borrow

import (
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/requests"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	borrowRepository "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrow_log"
	itemRepository "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	itemsetRepository "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item_set"
	logsystem "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/log"
	userRepository "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/user"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/item/factory"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Return(ctx context.Context, userID string, req *requests.ReturnRequest) error
	Borrow(ctx context.Context, userID string, itemID string) error

	GetUsersBorrowedItems(ctx context.Context, userID string) ([]responses.UserBorrrowResponse, error)
	GetAllBorrowedItems(ctx context.Context) ([]responses.AdminBorrowResponse, error)
}

type service struct {
	userRepo    userRepository.Repository
	borrowRepo  borrowRepository.Repository
	itemRepo    itemRepository.Repository
	itemSetRepo itemsetRepository.Repository
	logRepo     logsystem.Repository
}

func NewBorrowService(
	borrowRepo borrowRepository.Repository,
	itemRepo itemRepository.Repository,
	itemSetRepo itemsetRepository.Repository,
	logRepo logsystem.Repository) Service {
	return &service{
		borrowRepo:  borrowRepo,
		itemRepo:    itemRepo,
		itemSetRepo: itemSetRepo,
		logRepo:     logRepo,
	}
}

func (s *service) Borrow(ctx context.Context, userID string, itemID string) error {
	var borrowFactory factory.Borrowable

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return exceptions.ErrInvalidUUID
	}
	itemIDUUID, err := uuid.Parse(itemID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return exceptions.ErrInvalidUUID
	}

	item, err := s.itemRepo.GetItemByID(ctx, itemIDUUID)
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

	return borrowFactory.BorrowItem(ctx, userIDUUID, item, &children)

}

// Return implements Service.
func (s *service) Return(ctx context.Context, userID string, req *requests.ReturnRequest) error {
	var itemBorrowableFactory factory.Borrowable
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return exceptions.ErrInvalidUUID
	}
	borrowIDUUID, err := uuid.Parse(req.BorrowID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return exceptions.ErrInvalidUUID
	}

	borrow, err := s.borrowRepo.FindBorrowLogByUserIDAndBorrowID(ctx, userIDUUID, borrowIDUUID)
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

	borrow.ReturnImgURL = &req.ImageURL

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
	return itemBorrowableFactory.ReturnItem(ctx, borrow, &children)
}

// GetUsersBorrowedItems implements Service.
func (s *service) GetUsersBorrowedItems(ctx context.Context, userID string) ([]responses.UserBorrrowResponse, error) {
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Error().Err(err).Msg("invalid uuid format")
		return nil, exceptions.ErrInvalidUUID
	}

	borrows, err := s.borrowRepo.FindBorrowLogByUserID(ctx, userIDUUID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get borrow logs by user id")
		return nil, err
	}

	var results []responses.UserBorrrowResponse
	for _, borrow := range borrows {
		item, err := s.itemRepo.GetItemByID(ctx, borrow.ItemID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get item by id")
			return nil, err
		}
		if item == nil {
			log.Error().Err(err).Msg("item not found")
			return nil, exceptions.ErrItemNotFound
		}

		result := responses.UserBorrrowResponse{
			BorrowID:     borrow.BorrowID.String(),
			ItemName:     item.ItemName,
			BorrowDate:   utils.ToStringDateTime(borrow.BorrowDate),
			BorrowStatus: borrow.BorrowStatus,
		}

		if borrow.ReturnDate != nil {
			timeResult := utils.ToStringDateTime(*borrow.ReturnDate)
			result.ReturnDate = &timeResult
		}
		results = append(results, result)
	}

	return results, nil
}

// GetAllBorrowedItems implements Service.
func (s *service) GetAllBorrowedItems(ctx context.Context) ([]responses.AdminBorrowResponse, error) {
	borrows, err := s.borrowRepo.GetAllBorrowLogs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all borrow logs")
		return nil, err
	}

	var results []responses.AdminBorrowResponse
	for _, borrow := range borrows {
		user, err := s.userRepo.FindByID(ctx, borrow.UserID.String())
		if err != nil {
			log.Error().Err(err).Msg("failed to get user by id")
			return nil, err
		}
		item, err := s.itemRepo.GetItemByID(ctx, borrow.ItemID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get item by id")
			return nil, err
		}
		if item == nil {
			log.Error().Err(err).Msg("item not found")
			return nil, exceptions.ErrItemNotFound
		}

		result := responses.AdminBorrowResponse{
			BorrowID:     borrow.BorrowID.String(),
			UserName:     user.UserFullName,
			ItemName:     item.ItemName,
			BorrowDate:   utils.ToStringDateTime(borrow.BorrowDate),
			UserID:       user.UserID.String(),
			ItemID:       item.ItemID.String(),
			BorrowStatus: borrow.BorrowStatus,
		}

		if borrow.BorrowParentID != nil {
			parentID := borrow.BorrowParentID.String()
			result.BorrowParentID = &parentID
		}

		if borrow.ReturnDate != nil {
			timeResult := utils.ToStringDateTime(*borrow.ReturnDate)
			result.ReturnDate = &timeResult
		}

		if borrow.ReturnImgURL != nil {
			result.ReturnURL = borrow.ReturnImgURL
		}
		results = append(results, result)
	}

	return results, nil
}
