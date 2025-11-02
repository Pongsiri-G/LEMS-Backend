package state

// import "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"

type State interface {
	Return(ctx *BorrowStateContext) error
	Borrow(ctx *BorrowStateContext) error
}