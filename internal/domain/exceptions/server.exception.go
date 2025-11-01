package exceptions

import "errors"

var ErrNoTransaction = errors.New("no transaction in context")