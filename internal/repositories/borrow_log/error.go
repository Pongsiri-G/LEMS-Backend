package borrowlog

import "errors"

var (
	ErrMoreThanOneBorrowLog = errors.New("more than one borrow log found")
)
