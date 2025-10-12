package user

import "errors"

var (
	ErrNotFound             = errors.New("user not found")
	ErrRejectOnlyPending    = errors.New("only pending user can be rejected")
	ErrDeactivateNotPending = errors.New("cannot deactivate a PENDING account")
)
