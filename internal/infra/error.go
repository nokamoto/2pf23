package infra

import "fmt"

// ErrNotFound indicates that the resource does not exist.
var ErrNotFound = fmt.Errorf("not found")
