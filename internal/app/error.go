package app

import "errors"

// ErrInvalidArgument indicates client specified an invalid argument.
var ErrInvalidArgument = errors.New("invalid argument")

// ErrNotFound indicates some requested resource was not found.
var ErrNotFound = errors.New("not found")
