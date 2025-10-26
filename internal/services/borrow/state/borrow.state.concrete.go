package state

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
)

type BorrowedState struct{}
type ReturnedState struct{}
type WaitingState struct{}
type CanceledState struct{}

// ----------------------------------------------------------------------------------------------------
// BorrowedState

func (b *BorrowedState) Return(ctx *BorrowStateContext) error {
	now := utils.BangkokNow()
	ctx.state = &ReturnedState{}
	ctx.borrowLog.BorrowStatus = enums.StatusReturned
	ctx.borrowLog.ReturnDate = &now
	ctx.borrowLog.UpdatedAt = now
	err := ctx.borrowRepo.EditBorrowLog(ctx.ctx, ctx.borrowLog)
	if err != nil {
		log.Error().Err(err).Msg("failed to update borrow log")
		return err
	}
	return nil
}

func (b *BorrowedState) Borrow(ctx *BorrowStateContext) error {
	fmt.Println("this borrow log can't borrow")
	return nil
}

func (b *BorrowedState) Cancel(ctx *BorrowStateContext) error {
	fmt.Println("this borrow log can't cancel")
	return nil
}

// ----------------------------------------------------------------------------------------------------
// ReturnedState

func (b *ReturnedState) Return(ctx *BorrowStateContext) error {
	fmt.Println("this borrow log already returned")
	return nil
}

func (b *ReturnedState) Borrow(ctx *BorrowStateContext) error {
	fmt.Println("this borrow log can't borrow")
	return nil
}

func (b *ReturnedState) Cancel(ctx *BorrowStateContext) error {
	fmt.Println("this borrow log can't cancel")
	return nil
}

// ----------------------------------------------------------------------------------------------------
// WaitingState

func (b *WaitingState) Return(ctx *BorrowStateContext) error {
	fmt.Println("this borrow log can't return")
	return nil
}

func (b *WaitingState) Borrow(ctx *BorrowStateContext) error {
	// Implement logic for check if this Borrow is head of borrow queue then change state to BorrowState
	return nil
}

func (b *WaitingState) Cancel(ctx *BorrowStateContext) error {
	ctx.state = &CanceledState{}
	err := ctx.borrowRepo.EditBorrowLog(ctx.ctx, ctx.borrowLog)
	if err != nil {
		log.Error().Err(err).Msg("failed to update borrow log")
		return err
	}
	return nil
}

// ----------------------------------------------------------------------------------------------------
// CancelState

func (b *CanceledState) Return(ctx *BorrowStateContext) error {
	fmt.Println("this borrow log already cancel can't return")
	return nil
}

func (b *CanceledState) Borrow(ctx *BorrowStateContext) error {
	fmt.Println("this borrow log already cancel can't borrow")
	return nil
}

func (b *CanceledState) Cancel(ctx *BorrowStateContext) error {
	fmt.Println("this borrow log already cancel can't return")
	return nil
}
