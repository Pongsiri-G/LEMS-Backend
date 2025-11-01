package exceptions

import "errors"

var ErrNoTransaction = errors.New("no transaction in context")
var ErrNotImplemented = errors.New("feature not implemented")
var ErrInvalidExportType = errors.New("invalid export type")
