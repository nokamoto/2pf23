package helper

import (
	"testing"

	"github.com/bufbuild/connect-go"
	v1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
)

func TestConnectEqual(t *testing.T) {
	x := &v1alpha.Cluster{
		Name: "test",
	}
	m := ConnectEqual(x)

	if !m.Matches(connect.NewRequest(&v1alpha.Cluster{Name: "test"})) {
		t.Errorf("should be equal")
	}
	if m.Matches(connect.NewRequest(&v1alpha.Cluster{Name: "test2"})) {
		t.Errorf("should not be equal")
	}
	if m.Matches(nil) {
		t.Errorf("should not be equal")
	}
}
