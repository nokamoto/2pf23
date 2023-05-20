package helper

import (
	"testing"

	v1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
)

func TestProtoEqual(t *testing.T) {
	x := &v1alpha.Cluster{
		Name: "test",
	}
	m := ProtoEqual(x)

	if !m.Matches(&v1alpha.Cluster{Name: "test"}) {
		t.Errorf("should be equal")
	}
	if m.Matches(&v1alpha.Cluster{Name: "test2"}) {
		t.Errorf("should not be equal")
	}
	if m.Matches(nil) {
		t.Errorf("should not be equal")
	}
}
