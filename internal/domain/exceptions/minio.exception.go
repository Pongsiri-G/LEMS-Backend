package exceptions

import "errors"

var (
	ErrInvalidFileFormat = errors.New("invalid file format")
)
