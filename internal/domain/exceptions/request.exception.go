package exceptions

import "errors"

var (
	ErrRequestNotFound      = errors.New("request not found")
	ErrRequestItemInvalid   = errors.New("request item invalid")
	ErrRequestItemIDInvalid = errors.New("request item ID invalid")
)
