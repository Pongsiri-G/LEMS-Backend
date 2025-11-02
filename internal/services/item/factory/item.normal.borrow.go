package factory

import (
	"context"
	"sync"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	BorrowRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrow_log"
	ItemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	logsystem "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/log"
	stateBorrowLog "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/borrow/state"
	stateItem "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/item/state"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var normalBorrowableLock = sync.Mutex{}
var normalBorrowable Borrowable = nil

type ItemBorrowable struct {
	itemRepo   ItemRepo.Repository
	borrowRepo BorrowRepo.Repository
	logRepo    logsystem.Repository
}

func NewNormalItemBorrowable(
	itemRepo ItemRepo.Repository,
	borrowRepo BorrowRepo.Repository,
	logRepo logsystem.Repository,
) Borrowable {
	if normalBorrowable == nil {
		normalBorrowableLock.Lock()
		defer normalBorrowableLock.Unlock()
		if normalBorrowable == nil {
			normalBorrowable = &ItemBorrowable{
				itemRepo:   itemRepo,
				borrowRepo: borrowRepo,
				logRepo:    logRepo,
			}
		}
	}
	return normalBorrowable
}

// BorrowItem implements Borrowable.
func (i *ItemBorrowable) BorrowItem(ctx context.Context, userID uuid.UUID, item *models.Item, children *[]models.Item) error {
	borrowLog := models.BorrowLog{
		BorrowID:     uuid.New(),
		UserID:       userID,
		ItemID:       item.ItemID,
		BorrowStatus: enums.StatusBorrowed,
		BorrowDate:   utils.BangkokNow(),
		ReturnDate:   nil,
		CreatedAt:    utils.BangkokNow(),
		UpdatedAt:    utils.BangkokNow(),
	}
	itemContext := stateItem.NewStateContext(ctx, *item, i.itemRepo)

	// if item.ItemCurrentQuantity-1 < 0 {
	// 	log.Error().Err(exceptions.ErrItemQuantityInSufficient).Msg("quantity of the item less than zero")
	// 	return exceptions.ErrItemQuantityInSufficient
	// }

	// if item.ItemCurrentQuantity-1 == 0 {
	// 	item.ItemStatus = enums.ItemStatusOutOfStock
	// }

	// item.ItemCurrentQuantity -= 1
	// err := i.itemRepo.UpdateItem(ctx, item)
	// if err != nil {
	// 	log.Error().Err(err).Msg("failed to update quantity of item")
	// 	return exceptions.ErrFailedToUpdateQuantity
	// }
	err := itemContext.GetState().Borrow(itemContext)
	if err != nil {
		return err
	}

	err = i.borrowRepo.CreateBorrowLog(ctx, borrowLog)
	if err != nil {
		log.Error().Err(err).Msg("failed to create borrow log")
		return err
	}

	return nil
}

// ReturnItem implements Borrowable.
func (i *ItemBorrowable) ReturnItem(ctx context.Context, borrowLog *models.BorrowLog, children *[]models.BorrowLog) error {
	log.Info().Msg("Returning normal item")
	borrowLogContext := stateBorrowLog.NewStateContext(ctx, *borrowLog, i.borrowRepo)
	item, err := i.itemRepo.GetItemByID(ctx, borrowLog.ItemID)
	if err != nil {
		return err
	}
	if item == nil {
		log.Error().Err(err).Msg("item not found")
		return exceptions.ErrItemNotFound
	}
	itemContext := stateItem.NewStateContext(ctx, *item, i.itemRepo)

	err = borrowLogContext.GetState().Return(borrowLogContext)
	if err != nil {
		log.Error().Err(err).Msg("failed to update borrow log")
		return err
	}
	// borrowLog.BorrowStatus = enums.StatusReturned
	// now := utils.BangkokNow()
	// borrowLog.ReturnDate = &now
	// borrowLog.UpdatedAt = now

	// err = i.borrowRepo.EditBorrowLog(ctx, borrowLog)
	// if err != nil {
	// 	log.Error().Err(err).Msg("failed to update borrow log")
	// 	return err
	// }

	err = itemContext.GetState().Return(itemContext)
	if err != nil {
		log.Error().Err(err).Msg("failed to update quantity of item")
		return err
	}
	// if item.ItemCurrentQuantity == 0 {
	// 	item.ItemStatus = enums.ItemStatusAvailable
	// }
	// item.ItemCurrentQuantity += 1
	// err = i.itemRepo.UpdateItem(ctx, item)
	// if err != nil {
	// 	log.Error().Err(err).Msg("failed to update quantity of item")
	// 	return err
	// }

	return nil
}
