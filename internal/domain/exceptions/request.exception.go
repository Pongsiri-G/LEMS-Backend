package exceptions

import "errors"

var (
	ErrRequestNotFound = errors.New("request not found")
)
