package state

import (
	// "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	borrowRepository "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/borrow_log"
)

type BorrowStateContext struct {
	ctx        context.Context
	borrowLog  *models.BorrowLog
	state      State
	borrowRepo borrowRepository.Repository
}

func NewStateContext(ctx context.Context, borrowLog models.BorrowLog, borrowRepo borrowRepository.Repository) *BorrowStateContext {
	var s State
	switch borrowLog.BorrowStatus {
	case enums.StatusBorrowed:
		s = &BorrowedState{}
	case enums.StatusReturned:
		s = &ReturnedState{}
	default:
		s = &BorrowedState{}
	}
	return &BorrowStateContext{
		ctx:        ctx,
		borrowLog:  &borrowLog,
		borrowRepo: borrowRepo,
		state:      s,
	}
}

func (b *BorrowStateContext) GetCtx() context.Context {
	return b.ctx
}

func (b *BorrowStateContext) GetBorrowLog() *models.BorrowLog {
	return b.borrowLog
}

func (b *BorrowStateContext) GetState() State {
	return b.state
}

func (b *BorrowStateContext) SetState(s State) {
	b.state = s
}

func (b *BorrowStateContext) GetBorrowRepo() borrowRepository.Repository {
	return b.borrowRepo
}
