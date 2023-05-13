package infra

import "fmt"

// ErrNotFound indicates that the resource does not exist.
var ErrNotFound = fmt.Errorf("not found")

// ErrInvalidArgument indicates that the argument is invalid.
var ErrInvalidArgument = fmt.Errorf("invalid argument")

// ErrAlreadyExists indicates that the resource already exists.
var ErrAlreadyExists = fmt.Errorf("already exists")
