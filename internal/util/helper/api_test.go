package helper

import (
	"errors"
	"testing"

	"github.com/bufbuild/connect-go"
	kev1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
)

func TestCodeOf(t *testing.T) {
	if got := CodeOf(nil); got != 0 {
		t.Errorf("CodeOf() = %v, want %v", got, 0)
	}
	if got := CodeOf(errors.New("")); got != connect.CodeUnknown {
		t.Errorf("CodeOf() = %v, want %v", got, connect.CodeUnknown)
	}
	if got := CodeOf(connect.NewError(connect.CodeInternal, errors.New(""))); got != connect.CodeInternal {
		t.Errorf("CodeOf() = %v, want %v", got, connect.CodeInternal)
	}
}

func TestGetMsg(t *testing.T) {
	req := connect.NewRequest(&kev1alpha.CreateClusterRequest{})
	if got := GetMsg(req); got != req.Msg {
		t.Errorf("GetMsg() = %v, want %v", got, req.Msg)
	}
	if got := GetMsg[kev1alpha.CreateClusterRequest](nil); got != nil {
		t.Errorf("GetMsg() = %v, want %v", got, nil)
	} else {
		// t.Logf("panic: %v", got.Cluster)
		t.Logf("should not panic: %v", got.GetCluster())
	}
}

func TestGetResponseMsg(t *testing.T) {
	res := connect.NewResponse(&kev1alpha.CreateClusterRequest{})
	if got := GetResponseMsg(res); got != res.Msg {
		t.Errorf("GetMsg() = %v, want %v", got, res.Msg)
	}
	if got := GetResponseMsg[kev1alpha.CreateClusterRequest](nil); got != nil {
		t.Errorf("GetMsg() = %v, want %v", got, nil)
	} else {
		// t.Logf("panic: %v", got.Cluster)
		t.Logf("should not panic: %v", got.GetCluster())
	}
}
