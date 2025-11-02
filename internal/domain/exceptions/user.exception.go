package exceptions

import "errors"

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrNotAllowAccess = errors.New("you can't access to this.")
var ErrUserNotFound = errors.New("user not found")
var ErrInactiveUser = errors.New("user has been created\nbut user is not activated")