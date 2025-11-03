package state

import (
	"fmt"

	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/enums"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/domain/exceptions"
	"github.com/471-68-SE-Classroom/p1-final-project-backend-lems-ya/internal/utils"
	"github.com/rs/zerolog/log"
)

type AvailableState struct{}
type UnavailableState struct{}
type InLabOnlyState struct{}

// ----------------------------------------------------------------------------------------------------
// AvailableState

func (a *AvailableState) Return(ctx *ItemStateContext) error {
	now := utils.BangkokNow()
	if ctx.item.ItemCurrentQuantity == 0 {
		ctx.item.ItemStatus = enums.ItemStatusAvailable
	}
	ctx.item.ItemUpdatedAt = now
	ctx.item.ItemCurrentQuantity += 1
	err := ctx.itemRepo.UpdateItem(ctx.ctx, ctx.item)
	if err != nil {
		log.Error().Err(err).Msg("failed to update quantity of item")
		return err
	}
	return nil
}

func (a *AvailableState) Borrow(ctx *ItemStateContext) error {
	now := utils.BangkokNow()
	if ctx.item.ItemCurrentQuantity-1 < 0 {
		log.Error().Err(exceptions.ErrItemQuantityInSufficient).Msg("quantity of the item less than zero")
		return exceptions.ErrItemQuantityInSufficient
	}

	if ctx.item.ItemCurrentQuantity-1 == 0 {
		ctx.state = &UnavailableState{}
		ctx.item.ItemStatus = enums.ItemStatusOutOfStock
	}
	ctx.item.ItemUpdatedAt = now

	ctx.item.ItemCurrentQuantity -= 1
	err := ctx.itemRepo.UpdateItem(ctx.ctx, ctx.item)
	if err != nil {
		log.Error().Err(err).Msg("failed to update quantity of item")
		return exceptions.ErrFailedToUpdateQuantity
	}
	return nil
}

// ----------------------------------------------------------------------------------------------------
// UnavailableState

func (b *UnavailableState) Return(ctx *ItemStateContext) error {
	now := utils.BangkokNow()
	if ctx.item.ItemCurrentQuantity == 0 {
		ctx.state = &AvailableState{}
		ctx.item.ItemStatus = enums.ItemStatusAvailable
	}
	ctx.item.ItemCurrentQuantity += 1
	ctx.item.ItemUpdatedAt = now
	err := ctx.itemRepo.UpdateItem(ctx.ctx, ctx.item)
	if err != nil {
		log.Error().Err(err).Msg("failed to update quantity of item")
		return err
	}
	return nil
}

func (b *UnavailableState) Borrow(ctx *ItemStateContext) error {
	fmt.Println("this item currently unavailable state can't borrow")
	return nil
}

// ----------------------------------------------------------------------------------------------------
// InLabOnlyState

func (i *InLabOnlyState) Return(ctx *ItemStateContext) error {
	fmt.Println("this item currently in-lab only state can't borrow")
	return nil
}

func (i *InLabOnlyState) Borrow(ctx *ItemStateContext) error {
	fmt.Println("this item currently in-lab only state can't return")
	return nil
}
