package exceptions

import "errors"

var (
	ErrInvalidS3Url                 = errors.New("invalid s3 url")
	ErrInternalServer               = errors.New("internal server error")
	ErrFailedToParseAdminID         = errors.New("failed to parse admin ID")
	ErrFailedToParseTargetID        = errors.New("failed to parse target user ID")
	ErrFailedToCreateAdminActionLog = errors.New("failed to create admin action log")
)
