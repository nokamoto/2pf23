package cli

import (
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
)

// Merge merges the given packages into a single package.
func Merge(packages ...*v1.Package) (*v1.Package, error) {
	if len(packages) == 0 {
		return nil, nil
	}
	panic("not implemented")
}
