package exceptions

import "errors"

var ErrNoSuchStrategy = errors.New("the specified strategy does not exist")