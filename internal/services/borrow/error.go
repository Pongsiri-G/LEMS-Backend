package borrow

import "errors"

var (
	ErrInvalidUUID              = errors.New("invalid UUID")
	ErrItemQuantityInSufficient = errors.New("the quantity of the item is insufficient to allow borrowing")
	ErrFailedToUpdateQuantity   = errors.New("failed to update quantity of the item")
	ErrItemNotFound             = errors.New("item not found")
)
