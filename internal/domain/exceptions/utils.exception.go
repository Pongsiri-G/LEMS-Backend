package exceptions

import "errors"

var (
	ErrInvalidS3Url   = errors.New("invalid s3 url")
	ErrInternalServer = errors.New("internal server error")
)
