package state

// import "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"

type State interface {
	Borrow(ctx *ItemStateContext) error
	Return(ctx *ItemStateContext) error
}