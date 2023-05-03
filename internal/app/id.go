package app

import "github.com/google/uuid"

// ResourceIDGenerator generates a new resource ID.
type ResourceIDGenerator struct{}

func (r *ResourceIDGenerator) NewID() string {
	return uuid.New().String()
}
