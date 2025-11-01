package exceptions

import "errors"

var (
	ErrTagNotFound        = errors.New("tag not found")
	ErrTagAlreadyExists   = errors.New("tag already exists")
	ErrTagAlreadyAssigned = errors.New("tag already assigned to item")
	ErrTagNotAssigned     = errors.New("tag not assigned to item")
)
