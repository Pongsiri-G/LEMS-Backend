package exceptions

import "errors"

var (
	ErrItemQuantityInSufficient      = errors.New("the quantity of the item is insufficient to allow borrowing")
	ErrFailedToUpdateQuantity        = errors.New("failed to update quantity of the item")
	ErrBorrowLogNotFound             = errors.New("borrow log not found")
	ErrCannotReturnChildItemDirectly = errors.New("cannot return child item directly")
)
