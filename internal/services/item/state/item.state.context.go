package state

import (
	// "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"context"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/models"
	itemRepository "github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/repositories/item"
)

type ItemStateContext struct {
	ctx        context.Context
	item  *models.Item
	state      State
	itemRepo itemRepository.Repository
}

func NewStateContext(ctx context.Context, item models.Item, itemRepo itemRepository.Repository) *ItemStateContext {
	var s State
	switch item.ItemStatus {
	case enums.ItemStatusAvailable:
		s = &AvailableState{}
	case enums.ItemStatusOutOfStock:
		s = &UnavailableState{}
	default:
		s = &AvailableState{}
	}
	return &ItemStateContext{
		ctx:        ctx,
		item:  &item,
		itemRepo: itemRepo,
		state:      s,
	}
}

func (i *ItemStateContext) GetCtx() context.Context {
	return i.ctx
}

func (i *ItemStateContext) GetItem() *models.Item {
	return i.item
}

func (i *ItemStateContext) GetState() State {
	return i.state
}

func (i *ItemStateContext) SetState(s State) {
	i.state = s
}

func (i *ItemStateContext) GetBorrowRepo() itemRepository.Repository {
	return i.itemRepo
}
