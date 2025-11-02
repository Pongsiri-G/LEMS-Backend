package exceptions

import "errors"

var (
	ErrRequestNotFound               = errors.New("request not found")
	ErrRequestItemInvalid            = errors.New("request item invalid")
	ErrRequestItemIDInvalid          = errors.New("request item ID invalid")
	ErrRequestInvalidRequestType     = errors.New("invalid request type")
	ErrRequestNotExpectItemID        = errors.New("request type does not expect item ID")
	ErrRequestNotExpectItem          = errors.New("request type does not expect item")
	ErrUserIDIsNil                   = errors.New("user ID is nil")
	ErrRequestNotSupportedExportType = errors.New("request export type is not supported")
)
