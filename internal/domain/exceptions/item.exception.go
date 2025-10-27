package exceptions

import "errors"

var (
	ErrNoSuchStrategy        = errors.New("the specified strategy does not exist")
	ErrItemNotFound          = errors.New("item not found")
	ErrInvalidUUID           = errors.New("invalid uuid format")
	ErrRequestedItemNotFound = errors.New("requested item not found")
)
