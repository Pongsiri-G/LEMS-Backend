package factory

import (
	"context"
	"sync"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/responses"
	BorrowRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrow_log"
	ItemRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
	ItemSetRepo "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item_set"
	logsystem "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/log"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/services/borrow/state"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var childBorrowableLock = sync.Mutex{}
var childBorrowable Borrowable = nil

type ItemChildBorrowable struct {
	itemRepo    ItemRepo.Repository
	itemSetRepo ItemSetRepo.Repository
	borrowRepo  BorrowRepo.Repository
	logRepo     logsystem.Repository
}

func NewChildItemBorrowable(
	itemRepo ItemRepo.Repository,
	borrowRepo BorrowRepo.Repository,
	itemSetRepo ItemSetRepo.Repository,
	logRepo logsystem.Repository,
) Borrowable {
	if childBorrowable == nil {
		childBorrowableLock.Lock()
		defer childBorrowableLock.Unlock()
		if childBorrowable == nil {
			childBorrowable = &ItemChildBorrowable{
				itemRepo:    itemRepo,
				borrowRepo:  borrowRepo,
				itemSetRepo: itemSetRepo,
				logRepo:     logRepo,
			}
		}
	}
	return childBorrowable
}

// BorrowItem implements Borrowable.
func (i *ItemChildBorrowable) BorrowItem(ctx context.Context, userID uuid.UUID, item *models.Item, children *[]models.Item) error {
	resFactory := NewItemResponseFactoryWithChildrenConcrete(i.itemRepo, i.itemSetRepo)
	itemResp, err := resFactory.ToResponse(ctx, item, children)
	if err != nil {
		return err
	}

	parentBorrowLog := models.BorrowLog{
		BorrowID:     uuid.New(),
		UserID:       userID,
		ItemID:       item.ItemID,
		BorrowStatus: enums.StatusBorrowed,
		BorrowDate:   utils.BangkokNow(),
		ReturnDate:   nil,
		CreatedAt:    utils.BangkokNow(),
		UpdatedAt:    utils.BangkokNow(),
	}

	if item.ItemCurrentQuantity-1 < 0 {
		log.Error().Err(exceptions.ErrItemQuantityInSufficient).Msg("quantity of the item less than zero")
		return exceptions.ErrItemQuantityInSufficient
	}

	item.ItemCurrentQuantity -= 1
	borrowLogs, err := i.borrowChildItems(userID, parentBorrowLog.BorrowID, itemResp)
	if err != nil {
		log.Error().Err(err).Msg("failed to borrow child items")
		return err
	}

	err = i.borrowRepo.CreateBorrowLog(ctx, parentBorrowLog)
	if err != nil {
		log.Error().Err(err).Msg("failed to create borrow log for parent item")
		return err
	}

	if err := i.logRepo.CreateBorrowLog(ctx, userID, parentBorrowLog.ItemID); err != nil {
		log.Error().Err(err).Msg("failed to create log system borrow log")
		return err
	}

	for _, borrowLog := range borrowLogs {
		err = i.borrowRepo.CreateBorrowLog(ctx, borrowLog)
		if err != nil {
			log.Error().Err(err).Msg("failed to create borrow log for child item")
			return err
		}

		if err := i.logRepo.CreateBorrowLog(ctx, userID, borrowLog.ItemID); err != nil {
			log.Error().Err(err).Msg("failed to create log system borrow log for child item")
			return err
		}
	}

	err = i.itemRepo.UpdateItem(ctx, item)
	if err != nil {
		log.Error().Err(err).Msg("failed to update quantity of item")
		return err
	}

	for _, child := range *itemResp.Prerequisites {
		err = i.itemRepo.UpdateItem(ctx, &models.Item{
			ItemID:              child.ID,
			ItemQuantity:        child.Quantity,
			ItemName:            child.Name,
			ItemDescription:     child.Description,
			ItemPictureURL:      child.PictureURL,
			ItemStatus:          child.Status,
			ItemCurrentQuantity: child.CurrentQuantity - 1,
			ItemCreatedAt:       child.CreatedAt,
			ItemUpdatedAt:       utils.BangkokNow(),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to update quantity of child item")
			return err
		}
	}

	return nil
}

func (i *ItemChildBorrowable) borrowChildItems(userID uuid.UUID, parentID uuid.UUID, itemResp *responses.ItemResponse) ([]models.BorrowLog, error) {
	var childrenBorrowLogs []models.BorrowLog
	for _, child := range *itemResp.Prerequisites {

		childBorrowLog := models.BorrowLog{
			BorrowID:       uuid.New(),
			BorrowParentID: &parentID,
			UserID:         userID,
			ItemID:         child.ID,
			BorrowStatus:   enums.StatusBorrowed,
			BorrowDate:     utils.BangkokNow(),
			ReturnDate:     nil,
			CreatedAt:      utils.BangkokNow(),
			UpdatedAt:      utils.BangkokNow(),
		}

		if child.CurrentQuantity-1 < 0 {
			return nil, exceptions.ErrItemQuantityInSufficient
		}

		childrenBorrowLogs = append(childrenBorrowLogs, childBorrowLog)

	}
	return childrenBorrowLogs, nil
}

type allItemStruct struct {
	Item      *models.Item
	BorrowLog *models.BorrowLog
}

// ReturnItem implements Borrowable.
func (i *ItemChildBorrowable) ReturnItem(ctx context.Context, borrowLog *models.BorrowLog, children *[]models.BorrowLog) error {
	log.Info().Msg("Returning child item")
	borrowLogContext := state.NewStateContext(ctx, *borrowLog, i.borrowRepo)
	now := utils.BangkokNow()
	var allItem []allItemStruct
	item, err := i.itemRepo.GetItemByID(ctx, borrowLog.ItemID)
	if err != nil {
		return err
	}
	if item == nil {
		log.Error().Err(err).Msg("item not found")
		return exceptions.ErrItemNotFound
	}

	borrowLogContext.GetState().Return(borrowLogContext)
	// ---------------------------OLD CODE---------------------------
	// borrowLog.BorrowStatus = enums.StatusReturned
	// borrowLog.ReturnDate = &now
	// borrowLog.UpdatedAt = now
	item.ItemQuantity += 1
	item.ItemUpdatedAt = now
	allItem = append(allItem, allItemStruct{
		Item:      item,
		BorrowLog: borrowLog,
	})

	for _, childBorrowLog := range *children {
		childItem, err := i.itemRepo.GetItemByID(ctx, childBorrowLog.ItemID)
		if err != nil {
			return err
		}
		if childItem == nil {
			log.Error().Err(err).Msg("item not found")
			return exceptions.ErrItemNotFound
		}

		childBorrowLogContext := state.NewStateContext(ctx, childBorrowLog, i.borrowRepo)
		childBorrowLogContext.GetState().Return(childBorrowLogContext)

		// childBorrowLog.BorrowStatus = enums.StatusReturned
		// childBorrowLog.ReturnDate = &now
		// childBorrowLog.UpdatedAt = now
		childBorrowLog.ReturnImgURL = borrowLog.ReturnImgURL
		childItem.ItemCurrentQuantity += 1
		childItem.ItemUpdatedAt = now
		allItem = append(allItem, allItemStruct{
			Item:      childItem,
			BorrowLog: &childBorrowLog,
		})
	}

	for _, itemStruct := range allItem {
		// err = i.borrowRepo.EditBorrowLog(ctx, itemStruct.BorrowLog)
		// if err != nil {
		// 	log.Error().Err(err).Msg("failed to update borrow log")
		// 	return err
		// }

		err = i.itemRepo.UpdateItem(ctx, itemStruct.Item)
		if err != nil {
			log.Error().Err(err).Msg("failed to update quantity of item")
			return err
		}

		err = i.logRepo.CreateReturnLog(ctx, itemStruct.BorrowLog.UserID, itemStruct.BorrowLog.ItemID)
		if err != nil {
			log.Error().Err(err).Msg("failed to create log system return log")
			return err
		}
	}

	return nil

}
