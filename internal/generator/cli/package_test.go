package cli

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/nokamoto/2pf23/pkg/api/inhouse/v1"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestMerge(t *testing.T) {
	testcases := []struct {
		name     string
		packages []*v1.Package
		expected *v1.Package
	}{
		{
			name: "empty",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Merge(tc.packages...)
			if err != nil {
				t.Fatalf("failed to merge: %v", err)
			}
			if diff := cmp.Diff(tc.expected, actual, protocmp.Transform()); diff != "" {
				t.Errorf("Merge() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
