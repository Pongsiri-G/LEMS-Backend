package user

import "errors"

var (
	ErrNotFound          = errors.New("user not found")
	ErrUserIsNotPending  = errors.New("only pending user can be manage")
	ErrDeactivatePending = errors.New("cannot deactivate a PENDING account")
	ErrRevokeUser        = errors.New("cannot revoke a user")
	ErrAlreadyAdmin      = errors.New("user is already an admin")
)
